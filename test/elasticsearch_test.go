package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
)

const (
	esClientPodPrefix = "elasticsearch-kubeaddons-client-"
	esClientPort      = "9200"

	kibanaPodPrefix = "kibana-kubeaddons-"
	kibanaPort      = "5601"
)

func elasticsearchChecker(t *testing.T, cluster testcluster.Cluster) testharness.Job {
	return func(t *testing.T) error {
		localport, stop, err := portForwardPodWithPrefix(cluster, "kubeaddons", esClientPodPrefix, esClientPort)
		if err != nil {
			return fmt.Errorf("could not forward port to elasticsearch client pod: %s", err)
		}
		defer close(stop)

		if err := checkElasticsearchAvailable(localport); err != nil {
			return fmt.Errorf("failed to check elasticsearch is available: %s", err)
		}
		if err := checkElasticsearchCreateAndGetDoc(localport); err != nil {
			return fmt.Errorf("failed to create and get elasticsearch document: %s", err)
		}

		t.Logf("INFO: successfully tested elasticsearch")
		return nil
	}
}

func kibanaChecker(t *testing.T, cluster testcluster.Cluster) testharness.Job {
	return func(t *testing.T) error {
		localport, stop, err := portForwardPodWithPrefix(cluster, "kubeaddons", kibanaPodPrefix, kibanaPort)
		if err != nil {
			return fmt.Errorf("could not forward port to kibana pod: %s", err)
		}
		defer close(stop)

		// Check kibana status is healthy.
		// Kibana may serve 503s for some time after its pod is ready. Retry until kibana is ready to accept requests.
		var resp *http.Response
		path := "/api/status"
		retryWait := 10 * time.Second
		for tries := 0; tries < 20; tries++ {
			resp, err = http.Get(fmt.Sprintf("http://localhost:%d%s", localport, path))
			if err != nil {
				return fmt.Errorf("could not GET %s: %s", path, err)
			}

			if resp.StatusCode != http.StatusServiceUnavailable {
				break
			}
			time.Sleep(retryWait)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected GET %s status %d, got %d", path, http.StatusOK, resp.StatusCode)
		}

		b, err := ioutil.ReadAll(resp.Body)
		obj := map[string]interface{}{}
		if err := json.Unmarshal(b, &obj); err != nil {
			return fmt.Errorf("could not decode JSON response: %s", err)
		}

		status, ok := obj["status"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("JSON response missing key status with object value")
		}
		overallStatus, ok := status["overall"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("status missing key overall with object value")
		}
		overallState, ok := overallStatus["state"].(string)
		if !ok {
			return fmt.Errorf("overall status missing key state with string value")
		}
		if overallState != "green" {
			return fmt.Errorf("expected kibana state green, got %s", overallState)
		}

		t.Logf("INFO: successfully tested kibana")
		return nil
	}
}

// checkElasticsearchAvailable checks that the elasticsearch API is available on a local port.
// Returns `nil` if the API responds to `GET` `/` with JSON containing `cluster_uuid`, otherwise an error.
func checkElasticsearchAvailable(localport int) error {
	path := "/"
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d%s", localport, path))
	if err != nil {
		return fmt.Errorf("could not GET %s: %s", path, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected GET %s status %d, got %d", path, http.StatusOK, resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	obj := map[string]interface{}{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return fmt.Errorf("could not decode JSON response: %s", err)
	}

	key := "cluster_uuid"
	if _, ok := obj[key]; !ok {
		return fmt.Errorf("JSON response missing key %s", key)
	}

	return nil
}

// checkElasticsearchAvailable checks that the elasticsearch API on a local port can create and retrieve documents.
// Returns `nil` if a document can be created and then retrieved, otherwise an error.
func checkElasticsearchCreateAndGetDoc(localport int) error {
	docKey := "test_key"
	docVal := uuid.New().String()

	// Create document.
	doc := []byte(fmt.Sprintf(`{"%s": "%s"}`, docKey, docVal))

	path := "/test_index/_doc"
	resp, err := http.Post(fmt.Sprintf("http://localhost:%d%s", localport, path), "application/json", bytes.NewBuffer(doc))
	if err != nil {
		return fmt.Errorf("could not POST %s: %s", path, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected POST %s status %d, got %d", path, http.StatusCreated, resp.StatusCode)
	}

	// Get document ID from create response.
	b, err := ioutil.ReadAll(resp.Body)
	obj := map[string]interface{}{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return fmt.Errorf("could not decode JSON response: %s", err)
	}

	idKey := "_id"
	docID, ok := obj[idKey].(string)
	if !ok {
		return fmt.Errorf("JSON response missing key %s with string value", idKey)
	}

	// Get document by its ID.
	path = "/test_index/_doc/" + docID
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d%s", localport, path))
	if err != nil {
		return fmt.Errorf("could not GET %s: %s", path, err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected GET %s status %d, got %d", path, http.StatusOK, resp.StatusCode)
	}

	// Check document content.
	b, err = ioutil.ReadAll(resp.Body)
	obj = map[string]interface{}{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return fmt.Errorf("could not decode JSON response: %s", err)
	}

	sourceKey := "_source"
	docSource, ok := obj[sourceKey].(map[string]interface{})
	if !ok {
		return fmt.Errorf("JSON response missing key %s with object value", sourceKey)
	}

	val, ok := docSource[docKey].(string)
	if !ok {
		return fmt.Errorf("doc missing key %s with string value", docKey)
	}
	if val != docVal {
		return fmt.Errorf("expected doc key %s to have value %s, got %s", docKey, docVal, val)
	}

	return nil
}
