package test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	networkutils "github.com/mesosphere/ksphere-testing-framework/pkg/utils/networking"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// podReadyTimeout defines the maximum time to wait within portForwardPodWithPrefix for a pod and its containers
	// to be Running and Ready. If we hit the timeout, there are most likely issues with the addon itself.
	podReadyTimeout = 300 * time.Second
)

func portForwardPodWithPrefix(cluster testcluster.Cluster, ns, prefix, port string) (int, chan struct{}, error) {
	pod, err := findPodWithPrefix(cluster, ns, prefix)
	if err != nil {
		return 0, nil, fmt.Errorf("could not find pod with prefix %s: %s", prefix, err)
	}

	// Wait for Pod running and its containers to be ready for some time or timeout!
	timeout := time.After(podReadyTimeout)
	tick := time.Tick(500 * time.Millisecond)
	var containersReady bool
	for {
		select {
		case <-timeout:
			return 0, nil, fmt.Errorf("Timeout after %.0f minutes - Pod %s in phase %s, containers of it are overall not ready", podReadyTimeout.Minutes(), pod.Name, pod.Status.Phase)
		case <-tick:
			// Check first that Pod is in Running phase
			if pod.Status.Phase != corev1.PodRunning {
				continue
			}
			// Check second that containers are in READY phase
			for _, c := range pod.Status.ContainerStatuses {
				if !c.Ready && containersReady {
					containersReady = false
				} else {
					containersReady = true
				}
			}
			if !containersReady {
				continue
			}
			break
		}
		break
	}
	return networkutils.PortForward(cluster.Config(), pod.Namespace, pod.Name, port)
}

func logsFromPodWithPrefix(cluster testcluster.Cluster, ns, prefix string) (io.ReadCloser, error) {
	pod, err := findPodWithPrefix(cluster, ns, prefix)
	if err != nil {
		return nil, fmt.Errorf("could not find pod with prefix %s: %s", prefix, err)
	}

	logs, err := cluster.Client().CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{}).Stream(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("could not request logs for pod %s/%s: %s", pod.Namespace, pod.Name, err)
	}

	return logs, nil
}

func findPodWithPrefix(cluster testcluster.Cluster, ns, prefix string) (*corev1.Pod, error) {
	pods, err := cluster.Client().CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, pod := range pods.Items {
		if pod.ObjectMeta.GetDeletionTimestamp() != nil {
			continue
		}
		if strings.HasPrefix(pod.Name, prefix) {
			return &pod, nil
		}
	}

	return nil, fmt.Errorf("pod with name prefix %s in namespace %s not found", prefix, ns)
}
