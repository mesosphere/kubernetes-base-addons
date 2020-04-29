package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
)

const (
	promPodPrefix = "prometheus-prometheus-kubeaddons-prom-prometheus-"
	promPort      = "9090"

	alertmanagerPodPrefix = "alertmanager-prometheus-kubeaddons-prom-alertmanager-"
	alertmanagerPort      = "9093"

	grafanaPodPrefix = "prometheus-kubeaddons-grafana"
	grafanaPort      = "3000"
)

func promChecker(t *testing.T, cluster testcluster.Cluster) testharness.Job {
	return func(t *testing.T) error {
		localport, stop, err := portForwardPodWithPrefix(cluster, "kubeaddons", promPodPrefix, promPort)
		if err != nil {
			return fmt.Errorf("could not forward port to prometheus pod: %s", err)
		}
		defer close(stop)

		// Query prometheus and assert the response status is success.
		var resp *http.Response
		path := "/api/v1/labels"
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d%s", localport, path))
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

		status, ok := obj["status"].(string)
		if !ok {
			return fmt.Errorf("JSON response missing key status with string value")
		}
		if status != "success" {
			return fmt.Errorf("expected status success, got %s", status)
		}

		t.Logf("INFO: successfully tested prometheus")
		return nil
	}
}

func alertmanagerChecker(t *testing.T, cluster testcluster.Cluster) testharness.Job {
	return func(t *testing.T) error {
		localport, stop, err := portForwardPodWithPrefix(cluster, "kubeaddons", alertmanagerPodPrefix, alertmanagerPort)
		if err != nil {
			return fmt.Errorf("could not forward port to alertmanager pod: %s", err)
		}
		defer close(stop)

		// Check alertmanager status is ready.
		var resp *http.Response
		path := "/api/v2/status"
		resp, err = http.Get(fmt.Sprintf("http://localhost:%d%s", localport, path))
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

		clusterObj, ok := obj["cluster"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("JSON response missing key cluster with object value")
		}
		status, ok := clusterObj["status"].(string)
		if !ok {
			return fmt.Errorf("cluster missing key status with string value")
		}
		if status != "ready" {
			return fmt.Errorf("expected status ready, got %s", status)
		}

		t.Logf("INFO: successfully tested alertmanager")
		return nil
	}
}

func grafanaChecker(t *testing.T, cluster testcluster.Cluster) testharness.Job {
	return func(t *testing.T) error {
		localport, stop, err := portForwardPodWithPrefix(cluster, "kubeaddons", grafanaPodPrefix, grafanaPort)
		if err != nil {
			return fmt.Errorf("could not forward port to grafana pod: %s", err)
		}
		defer close(stop)

		// Check grafana is healthy.
		var resp *http.Response
		path := "/api/health"
		req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d%s", localport, path), nil)
		if err != nil {
			return fmt.Errorf("could not create GET %s: %s", path, err)
		}
		req.Header.Set("X-Forwarded-User", "admin")
		resp, err = (&http.Client{}).Do(req)
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

		status, ok := obj["database"].(string)
		if !ok {
			return fmt.Errorf("JSON response missing key database with string value")
		}
		if status != "ok" {
			return fmt.Errorf("expected status ok, got %s", status)
		}

		t.Logf("INFO: successfully tested grafana")
		return nil
	}
}
