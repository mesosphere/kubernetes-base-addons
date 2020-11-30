package newtest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mesosphere/dkp-test-framework/pkg/addons"
	"github.com/mesosphere/dkp-test-framework/pkg/cluster"
	"github.com/mesosphere/dkp-test-framework/pkg/cluster/kind"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

func deployKubeaddonsController(cluster cluster.Cluster) error {
	return addons.HelmInstall(cluster, "kubeaddons", "https://mesosphere.github.io/kubeaddons/chart", "kubeaddons", "kube-system", nil)
}

func deployReloader(restConfig rest.Config) error {
	restConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: "kubeaddons.mesosphere.io", Version: "v1beta1"}
	client, err := rest.RESTClientFor(&restConfig)
	if err != nil {
		return fmt.Errorf("couldn't create REST client: %w", err)
	}

	if err := client.Post().
		Resource("addons").
		Namespace("kube-system").
		Body("testdata/reloader.json").
		Do(context.Background()).Error(); err != nil {
		return fmt.Errorf("couldn't create reloader addon: %w", err)
	}

	return nil
}

func deployNginx(client kubernetes.Interface) (*corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "html",
		},
		Data: map[string]string{
			"index.html": `original`,
		},
	}
	_, err := client.CoreV1().ConfigMaps("default").Create(context.Background(), &cm, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not create nginx ConfigMap: %w", err)
	}

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
			Annotations: map[string]string{
				"reloader.stakater.com/auto": "true",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "nginx",
						Image: "nginx:1.19.4",
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "html",
							MountPath: "/usr/share/nginx/html",
							ReadOnly:  true,
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "html",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "html",
								},
							},
						},
					}},
				},
			},
		},
	}
	_, err = client.AppsV1().Deployments("default").Create(context.Background(), &deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not create nginx Deployment: %w", err)
	}

	if err := wait.Poll(5*time.Second, 60*time.Second, func() (bool, error) {
		deploy, err := client.AppsV1().Deployments("default").Get(context.Background(), "nginx", metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return deploy.Status.ReadyReplicas > 0, nil
	}); err != nil {
		return nil, fmt.Errorf("couldn't check deployment readiness: %w", err)
	}

	_, err = client.CoreV1().Services("default").Create(context.Background(), &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []corev1.ServicePort{{
				Port: 80,
			}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not create nginx Service: %w", err)
	}

	return &cm, nil
}

func curl(c typedcorev1.PodInterface, url string) (string, error) {
	pod, err := c.Create(context.Background(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "debug",
		},
		Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{{
				Name:  "debug",
				Image: "curlimages/curl:7.73.0",
				Args: []string{
					"-s",
					url,
				},
			}},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("could not create pod: %w", err)
	}

	defer func() {
		c.Delete(context.Background(), "debug", metav1.DeleteOptions{})
	}()

	if err := wait.Poll(1*time.Second, 60*time.Second, func() (bool, error) {
		var err error
		pod, err = c.Get(context.Background(), "debug", metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return pod.Status.Phase == corev1.PodSucceeded, nil
	}); err != nil {
		return "", fmt.Errorf("could not verify pod is running: %w", err)
	}

	res, err := c.GetLogs("debug", &corev1.PodLogOptions{}).Do(context.Background()).Raw()
	if err != nil {
		return "", fmt.Errorf("could not get pod logs: %w", err)
	}
	return string(res), nil
}

// This tests checks the proper functioning of the reloader Addon as defined in 'testdata/reloader.json' by
// conducting the following steps:
//
// 1. spin up kind cluster
// 2. deploy kubeaddons controller
// 3. create reloader Addon resource
// 4. create a ConfigMap
// 5. deploy nginx consuming that ConfigMap for serving `/index.html`
// 6. check that nginx serves the ConfigMaps' data
// 7. change the ConfigMap's data
// 8. check that nginx now serves the updated data, proofing that reloader properly updated nginx
func TestReloaderAddon(t *testing.T) {
	prov := kind.NewProvisioner("reloader-test")
	timeoutCtx, cancelCreate := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancelCreate()
	wf := prov.CreateWorkflow(timeoutCtx)
	if err := wf.Start(); err != nil {
		t.Fatalf("couldn't start workflow: %s", err)
	}

	defer func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()
		del := prov.DeleteWorkflow(timeoutCtx)
		del.Start()
		del.Wait()
	}()

	if err := wf.Wait(); err != nil {
		t.Fatalf("error provisioning kind cluster: %v", err)
	}

	cluster, err := prov.Cluster()
	if err != nil {
		t.Fatalf("couldn't get cluster: %s", err)
	}
	if err := deployKubeaddonsController(cluster); err != nil {
		t.Fatalf("couldn't install kubeaddons: %s", err)
	}

	restConfig, err := cluster.ToRESTConfig()
	if err != nil {
		t.Fatalf("couldn't get REST config: %s", err)
	}

	if err := deployReloader(*restConfig); err != nil {
		t.Fatalf("couldn't deploy reloader: %s", err)
	}

	cm, err := deployNginx(cluster.Client())
	if err != nil {
		t.Fatalf("couldn't deploy nginx: %s", err)
	}

	res, err := curl(cluster.Client().CoreV1().Pods("default"), "http://nginx.default")
	if err != nil {
		t.Fatalf("couldn't check nginx service: %s", err)
	}
	if res != "original" {
		t.Fatalf("unexpected page content retrieved. Expected 'original', got '%s'", res)
	}

	cm.Data["index.html"] = "updated"
	if _, err := cluster.Client().CoreV1().ConfigMaps("default").Update(context.Background(), cm, metav1.UpdateOptions{}); err != nil {
		t.Fatalf("couldn't update ConfigMap: %s", err)
	}

	if err := wait.Poll(2*time.Second, 60*time.Second, func() (bool, error) {
		res, err := curl(cluster.Client().CoreV1().Pods("default"), "http://nginx.default")
		if err != nil {
			t.Fatalf("couldn't check nginx service: %s", err)
		}
		return res == "updated", nil
	}); err != nil {
		t.Fatalf("error waiting for nginx to return updated content: %s", err)
	}
}
