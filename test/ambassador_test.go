package test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	ambassadorv2 "github.com/datawire/ambassador/pkg/api/getambassador.io/v2"
	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta2"
	"github.com/mesosphere/kubeaddons/pkg/constants"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	kscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

const waitForAmbassador = time.Minute * 1

func ambassadorChecker(t *testing.T, cluster testcluster.Cluster) testharness.Job {
	return func(t *testing.T) error {
		app := "quote"
		ns := "default"

		// create a test application to ensure proper connectivity
		testDeployment := appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app,
				Namespace: ns,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: pointer.Int32Ptr(1),
				Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app": app}},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": app},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "backend",
							Image: "docker.io/datawire/quote:0.4.1",
							Ports: []corev1.ContainerPort{{
								Name:          "http",
								ContainerPort: 8080,
							}}}}}}}}
		deployment, err := cluster.Client().AppsV1().Deployments(ns).Create(context.TODO(), &testDeployment, metav1.CreateOptions{})
		if err != nil {
			t.Fatal(err)
		}

		// create the svc for the
		testService := corev1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app,
				Namespace: ns,
			},
			Spec: corev1.ServiceSpec{
				Ports: []corev1.ServicePort{{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(8080),
				}},
				Selector: map[string]string{
					"app": deployment.Name,
				}}}
		svc, err := cluster.Client().CoreV1().Services(ns).Create(context.TODO(), &testService, metav1.CreateOptions{})
		if err != nil {
			t.Fatal(err)
		}

		// create a dynamic client for ambassador's API
		scheme := k8sruntime.NewScheme()
		if err := kscheme.AddToScheme(scheme); err != nil {
			t.Fatal(err)
		}
		if err := ambassadorv2.AddToScheme(scheme); err != nil {
			t.Fatal(err)
		}
		if err := v1beta2.AddToScheme(scheme); err != nil {
			t.Fatal(err)
		}
		mapper, err := apiutil.NewDynamicRESTMapper(cluster.Config())
		if err != nil {
			t.Fatal(err)
		}
		c, err := client.New(cluster.Config(), client.Options{Scheme: scheme, Mapper: mapper})
		if err != nil {
			t.Fatal(err)
		}

		// create a mapping for traffic to the service
		apptestMapping := ambassadorv2.Mapping{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-backend", app),
				Namespace: ns,
			},
			Spec: ambassadorv2.MappingSpec{
				Prefix:  "/backend/",
				Service: svc.Name,
			}}
		if err := c.Create(context.TODO(), &apptestMapping); err != nil {
			t.Fatal(err)
		}

		// I've checked with upstream, even though there's a status available in the Mapping API, they don't use it since several
		// versions ago, so for the time being we just give the mapping a reasonable amount of time to resolve.
		time.Sleep(time.Second * 10)

		// get the svc IP for ambassador
		localport, stop, err := portForwardPodWithPrefix(cluster, constants.DefaultAddonNamespace, "ambassador", "8080")
		if err != nil {
			return fmt.Errorf("could not forward port to elasticsearch client pod: %s", err)
		}
		defer close(stop)

		// make sure requests to the test application are successful
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/backend/", localport))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 OK from ambassador backend, got %s", resp.Status)
		}

		// check the contents of the response
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		data := make(map[string]string)
		if err := json.Unmarshal(b, &data); err != nil {
			t.Fatal(err)
		}
		qotd, ok := data["quote"]
		if !ok {
			t.Fatalf("structure of output from test app did not include \"quote\" key: %+v", data)
		}

		t.Logf("INFO: ambassador tests complete, mapping connectivity verified! quote of the day is: %s", qotd)
		return nil
	}
}
