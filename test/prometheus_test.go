package test

import (
	"fmt"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/mesosphere/kubeaddons/pkg/test"
)

func TestPrometheusGroup(t *testing.T) {
	if err := testgroup(t, "prometheus", prometheusChecker); err != nil {
		t.Fatal(err)
	}
}

func prometheusChecker(t *testing.T, cluster test.Cluster) test.Job {
	return func(t *testing.T) error {
		if err := kubectl("apply", "-f", "./artifacts/prometheus-checker.yaml"); err != nil {
			return err
		}

		succeeded := false
		timeout := time.Now().Add(time.Minute * 1)
		for timeout.After(time.Now()) {
			job, err := cluster.Client().BatchV1().Jobs("default").Get("prometheus-checker", metav1.GetOptions{})
			if err != nil {
				return err
			}
			if job.Status.Succeeded == 1 {
				succeeded = true
				break
			}
			time.Sleep(time.Second * 1)
		}

		if !succeeded {
			return fmt.Errorf("prometheus checker job did not succeed within timeout")
		}
		t.Log("prometheus checker job succeeded 🤪")
		return nil
	}
}
