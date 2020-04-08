package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
	networkutils "github.com/mesosphere/ksphere-testing-framework/pkg/utils/networking"
	corev1 "k8s.io/api/core/v1"
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
		pod, err := findPodWithPrefix(cluster, "kubeaddons", promPodPrefix)
		if err != nil {
			return fmt.Errorf("could not find prometheus pod: %s", err)
		}
		if pod.Status.Phase != corev1.PodRunning {
			return fmt.Errorf("prometheus pod %s is not running, it's in phase %s", pod.Name, pod.Status.Phase)
		}
		t.Logf("INFO: checking prometheus at pod/%s port %s", pod.Name, promPort)

		// Forward a local port to the prometheus API.
		localport, stop, err := networkutils.PortForward(cluster.Config(), pod.Namespace, pod.Name, promPort)
		if err != nil {
			return fmt.Errorf("could not set up port forward for pod/%s port %s: %s", pod.Name, promPort, err)
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
		pod, err := findPodWithPrefix(cluster, "kubeaddons", alertmanagerPodPrefix)
		if err != nil {
			return fmt.Errorf("could not find alertmanager pod: %s", err)
		}
		if pod.Status.Phase != corev1.PodRunning {
			return fmt.Errorf("alertmanager pod %s is not running, it's in phase %s", pod.Name, pod.Status.Phase)
		}
		t.Logf("INFO: checking alertmanager at pod/%s port %s", pod.Name, alertmanagerPort)

		// Forward a local port to the alertmanager API.
		localport, stop, err := networkutils.PortForward(cluster.Config(), pod.Namespace, pod.Name, alertmanagerPort)
		if err != nil {
			return fmt.Errorf("could not set up port forward for pod/%s port %s: %s", pod.Name, alertmanagerPort, err)
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
		pod, err := findPodWithPrefix(cluster, "kubeaddons", grafanaPodPrefix)
		if err != nil {
			return fmt.Errorf("could not find grafana pod: %s", err)
		}
		if pod.Status.Phase != corev1.PodRunning {
			return fmt.Errorf("grafana pod %s is not running, it's in phase %s", pod.Name, pod.Status.Phase)
		}
		t.Logf("INFO: checking grafana at pod/%s port %s", pod.Name, grafanaPort)

		// Forward a local port to the grafana API.
		localport, stop, err := networkutils.PortForward(cluster.Config(), pod.Namespace, pod.Name, grafanaPort)
		if err != nil {
			return fmt.Errorf("could not set up port forward for pod/%s port %s: %s", pod.Name, grafanaPort, err)
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
