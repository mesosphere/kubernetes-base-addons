package test

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	networkutils "github.com/mesosphere/ksphere-testing-framework/pkg/utils/networking"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	certManagerNamespace = "cert-manager"
	kubeConfigPath       = "/etc/kubernetes/admin.conf"
	rootCACertPath       = "/etc/kubernetes/pki/ca.crt"
	rootCAKeyPath        = "/etc/kubernetes/pki/ca.key"
)

func createCertManagerSecret(k testcluster.Cluster) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: certManagerNamespace,
		},
	}
	_, err := k.Client().CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("could not create cert-manager namespace: %w", err)
	}

	createSecretCommand := fmt.Sprintf(
		"set -o pipefail && "+
			"kubectl --kubeconfig %s create secret tls kubernetes-root-ca "+
			"--namespace=%s --cert=%s --key=%s --dry-run -o yaml "+
			"| kubectl --kubeconfig /etc/kubernetes/admin.conf apply -f -",
		kubeConfigPath,
		certManagerNamespace,
		rootCACertPath,
		rootCAKeyPath)
	a := []string{
		"-c",
		createSecretCommand,
	}
	cmd := exec.Command("bash", a...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running command in docker: %w", err)
	}
	return nil
}

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
