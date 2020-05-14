package test

import (
	"fmt"
	"strings"

	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	networkutils "github.com/mesosphere/ksphere-testing-framework/pkg/utils/networking"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func portForwardPodWithPrefix(cluster testcluster.Cluster, ns, prefix, port string) (int, chan struct{}, error) {
	pod, err := findPodWithPrefix(cluster, ns, prefix)
	if err != nil {
		return 0, nil, fmt.Errorf("could not find pod with prefix %s: %s", prefix, err)
	}
	if pod.Status.Phase != corev1.PodRunning {
		return 0, nil, fmt.Errorf("pod %s is not running, it's in phase %s", pod.Name, pod.Status.Phase)
	}
	return networkutils.PortForward(cluster.Config(), pod.Namespace, pod.Name, port)
}

func findPodWithPrefix(cluster testcluster.Cluster, ns, prefix string) (*corev1.Pod, error) {
	pods, err := cluster.Client().CoreV1().Pods(ns).List(metav1.ListOptions{})
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
