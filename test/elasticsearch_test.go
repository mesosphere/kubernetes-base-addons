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
		time.Sleep(time.Second * 120)
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
		time.Sleep(time.Second * 120)
		localport, stop, err := portForwardPodWithPrefix(cluster, "kubeaddons", kibanaPodPrefix, kibanaPort)
		if err != nil {
			return fmt.Errorf("could not forward port to kibana pod: %s", err)
		}
		defer close(stop)

		if err := waitForKibana(localport); err != nil {
			return fmt.Errorf("kibana took too long to become available: %s", err)
		}

		if err := checkKibanaStatus(localport); err != nil {
			return fmt.Errorf("failed to check kibana status: %s", err)
		}

		if err := checkKibanaDashboards(localport); err != nil {
			return fmt.Errorf("failed to check kibana dashboards: %s", err)
		}

		t.Logf("INFO: successfully tested kibana")
		return nil
	}
}

// waitForKibana returns nil once the Kibana API serves a 200 response, indicating it is ready to accept requests.
// Kibana may serve 503s for some time after its pod is ready.
func waitForKibana(localport int) error {
	path := "/api/status"
	retryWait := 10 * time.Second
	maxTries := 20

	var resp *http.Response
	var err error
	for tries := 0; tries < maxTries; tries++ {
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

	return nil
}

// checkKibanaStatus returns nil if the overall cluster state is green, otherwise an error.
func checkKibanaStatus(localport int) error {
	path := "/api/status"

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

	return nil
}

// checkKibanaDashboards returns nil if Kibana has all expected dashboards, otherwise an error.
func checkKibanaDashboards(localport int) error {
	expectedDashboards := []string{"Audit-Dashboard"}
	path := "/api/saved_objects/_find?type=dashboard"

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

	saved_objects, ok := obj["saved_objects"].([]interface{})
	if !ok {
		return fmt.Errorf("JSON response missing key saved_objects with array value")
	}

	titles := make(map[string]bool)
	for _, obj := range saved_objects {
		object, ok := obj.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Unexpected type for Kibana API object")
		}

		attributes, ok := object["attributes"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("Kibana API object missing key attributes with object value")
		}

		title, ok := attributes["title"].(string)
		if !ok {
			return fmt.Errorf("Kibana API object attributes missing key title with string value")
		}

		titles[title] = true
	}

	for _, title := range expectedDashboards {
		if _, ok := titles[title]; !ok {
			return fmt.Errorf("Dashboard %s not found", title)
		}
	}

	return nil
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
