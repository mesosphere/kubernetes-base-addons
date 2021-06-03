// +build nightly

package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/yaml.v2"

	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
)

const (
	updateFixturesEnvVar = "UPDATE_FIXTURES"

	kubeaddonsNamespace                     = "kubeaddons"
	prometheusMetricTestsFilename           = "testdata/prometheus-metric-tests.yaml"
	testNightlyGroupPromMetricTestsFilename = "testdata/test-nightly-group-prom-metric-tests.yaml"
)

func updateFixtures() bool {
	return true //os.Getenv(updateFixturesEnvVar) == "true"
}

func updateFixturesFatalMessage(err error) string {
	return fmt.Sprintf("failing, as "+updateFixturesEnvVar+"=true: %v", err)
}

func TestNightlyGroup(t *testing.T) {
	if err := testgroup(t, allAWSGroupName, defaultKindestNodeImage, nightlyPrometheusMetricsChecker); err != nil {
		t.Fatal(err)
	}
}

// TestUnmarshallPrometheusMetricNames is a quick unit test for unmarshallPrometheusMetrics
func TestUnmarshallPrometheusMetricNames(t *testing.T) {
	b, err := ioutil.ReadFile("testdata/prometheus-metric-output.json")
	if err != nil {
		t.Fatal(err)
	}
	metrics, err := unmarshallPrometheusMetrics(b)
	if err != nil {
		t.Fatal(err)
	}

	if updateFixtures() {
		commentHeader := "used in TestUnmarshallPrometheusMetricNames as unit test output of prometheus-metric-output.json"
		err := updatePrometheusMetricFixtures(prometheusMetricTestsFilename, commentHeader, metrics)
		if err != nil {
			t.Fatal(updateFixturesFatalMessage(err))
		}
	}
	err = testMetricFixtures(prometheusMetricTestsFilename, metrics)
	if err != nil {
		t.Fatal(err)
	}
}

func nightlyPrometheusMetricsChecker(t *testing.T, cluster testcluster.Cluster) testharness.Job {
	return func(t *testing.T) error {
		// let prometheus run for a bit
		time.Sleep(time.Second * 120)

		promMetricsBytes, err := getPrometheusMetricsBytes(cluster)
		if err != nil {
			return err
		}

		metrics, err := unmarshallPrometheusMetrics(promMetricsBytes)
		if err != nil {
			return err
		}

		if updateFixtures() {
			commentHeader := "used in TestNightlyGroup as e2e output of Konvoy in AWS with allAWSGroupName addons"
			err := updatePrometheusMetricFixtures(testNightlyGroupPromMetricTestsFilename, commentHeader, metrics)
			return fmt.Errorf(updateFixturesFatalMessage(err))
		}

		err = testMetricFixtures(prometheusMetricTestsFilename, metrics)
		if err != nil {
			return err
		}
		return nil
	}
}

func getPrometheusMetricsBytes(cluster testcluster.Cluster) ([]byte, error) {
	localport, stop, err := portForwardPodWithPrefix(cluster, "kubeaddons", promPodPrefix, promPort)
	if err != nil {
		return nil, fmt.Errorf("could not forward port to prometheus pod: %w", err)
	}
	defer close(stop)

	fullURL := prometheusMetricsFullURL(localport)
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("could not GET %s: %w", fullURL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected GET %s status %d, got %d", fullURL, http.StatusOK, resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func prometheusMetricsFullURL(localport int) string {
	queryURL := "/api/v1/query_range?"
	queryParams := strings.Join([]string{
		"query=up",
		// take a 5 min range to prevent flakes because data points are scattered over time
		"start=" + time.Now().Add(-5*time.Minute).Format(time.RFC3339),
		"end=" + time.Now().Format(time.RFC3339),
		"step=1s",
	}, "&")
	return fmt.Sprintf("http://localhost:%d%s", localport, queryURL+queryParams)
}

type prometheusUpMetricsResponse struct {
	Status string         `json:"status"`
	Data   prometheusData `json:"data"`
}

type prometheusData struct {
	ResultType string                        `json:"resultType"`
	Result     []prometheusUpMetricWithValue `json:"result"`
}

type prometheusUpMetricWithValue struct {
	Metric PrometheusUpMetric `json:"metric"`
}

type PrometheusUpMetric struct {
	Name      string `json:"__name__"`
	App       string `json:"app"`
	Namespace string `json:"namespace"`
	Service   string `json:"service"`
	Job       string `json:"job"`
}

func (p *PrometheusUpMetric) key() string {
	return strings.Join([]string{p.Name, p.App, p.Namespace, p.Service, p.Job}, " ; ")
}

func unmarshallPrometheusMetrics(data []byte) ([]PrometheusUpMetric, error) {
	resp := &prometheusUpMetricsResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, fmt.Errorf("error unmarshallPrometheusMetrics: %w", err)
	}

	if resp.Data.ResultType != "matrix" {
		return nil, fmt.Errorf("resultType of data is not matrix. got: %v", resp.Data.ResultType)
	}
	return summarizePrometheusMetrics(resp.Data.Result), nil
}

func summarizePrometheusMetrics(metricsWithValue []prometheusUpMetricWithValue) []PrometheusUpMetric {
	keyToMetric := map[string]PrometheusUpMetric{}

	for _, metricWithValue := range metricsWithValue {
		metric := metricWithValue.Metric
		keyToMetric[metric.key()] = metric
	}

	metrics := make([]PrometheusUpMetric, len(keyToMetric))
	i := 0
	for _, metric := range keyToMetric {
		metrics[i] = metric
		i++
	}
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].key() < metrics[j].key()
	})

	return metrics
}

func updatePrometheusMetricFixtures(filename, commentHeader string, expectedPromMetrics []PrometheusUpMetric) error {
	b, err := yaml.Marshal(expectedPromMetrics)
	if err != nil {
		return err
	}

	stringOutput := "# " + commentHeader + "\n" + string(b)

	fmt.Println(filename + " is now updated with:")
	fmt.Println(stringOutput)
	return ioutil.WriteFile(filename, []byte(stringOutput), 0600)
}

func testMetricFixtures(filename string, gotPromMetricsTests []PrometheusUpMetric) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var expectedPromMetricsTests []PrometheusUpMetric
	err = yaml.Unmarshal(b, &expectedPromMetricsTests)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(gotPromMetricsTests, expectedPromMetricsTests) {
		return fmt.Errorf(cmp.Diff(gotPromMetricsTests, expectedPromMetricsTests))
	}
	return nil
}
