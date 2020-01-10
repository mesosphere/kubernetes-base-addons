package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/blang/semver"
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"

	"sigs.k8s.io/kind/pkg/container/cri"
	volumetypes "github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"

	"sigs.k8s.io/kind/pkg/cluster/create"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha3"

	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

const (
	defaultKubernetesVersion = "1.15.6"
	patchStorageClass        = `{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}`
)

var addonTestingGroups = make(map[string][]string)

func init() {
	b, err := ioutil.ReadFile("groups.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, addonTestingGroups); err != nil {
		panic(err)
	}
}

func TestValidateUnhandledAddons(t *testing.T) {
	unhandled, err := findUnhandled()
	if err != nil {
		t.Fatal(err)
	}

	if len(unhandled) != 0 {
		names := make([]string, len(unhandled))
		for _, addon := range unhandled {
			names = append(names, addon.GetName())
		}
		t.Fatal(fmt.Errorf("the following addons are not handled as part of a testing group: %+v", names))
	}
}

func TestGeneralGroup(t *testing.T) {
	if err := testgroup(t, "general"); err != nil {
		t.Fatal(err)
	}
}

func TestElasticSearchGroup(t *testing.T) {
	if err := testgroup(t, "elasticsearch"); err != nil {
		t.Fatal(err)
	}
}

func TestPrometheusGroup(t *testing.T) {
	if err := testgroup(t, "prometheus"); err != nil {
		t.Fatal(err)
	}
}

func TestKommanderGroup(t *testing.T) {
	if err := testgroup(t, "kommander"); err != nil {
		t.Fatal(err)
	}
}

func TestIstioGroup(t *testing.T) {
	if err := testgroup(t, "istio"); err != nil {
		t.Fatal(err)
	}
}

func TestLocalVolumeProvisionerGroup(t *testing.T) {
	if err := testgroup(t, "localvolumeprovisioner"); err != nil {
		t.Fatal(err)
	}
}

func TestDispatchGroup(t *testing.T) {
	if err := testgroup(t, "dispatch"); err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Private Functions
// -----------------------------------------------------------------------------

func createNodeVolumes(numberVolumes int, nodePrefix string, node *v1alpha3.Node) error {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return fmt.Errorf("creating docker client: %w", err)
	}
	dockerClient.NegotiateAPIVersion(context.TODO())

	for index := 0; index < numberVolumes; index++ {
		volumeName := fmt.Sprintf("%s-%d", nodePrefix, index)

		volume, err := dockerClient.VolumeCreate(context.TODO(), volumetypes.VolumeCreateBody{
			Driver: "local",
			Name:   volumeName,
		})
		if err != nil {
			return fmt.Errorf("creating volume for node: %w", err)
		}

		node.ExtraMounts = append(node.ExtraMounts, cri.Mount{
			ContainerPath: fmt.Sprintf("/mnt/disks/%s", volumeName),
			HostPath:      volume.Mountpoint,
		})
	}

	return nil
}

func cleanupNodeVolumes(numberVolumes int, nodePrefix string, node *v1alpha3.Node) error {
	dockerClient, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		return fmt.Errorf("creating docker client: %w", err)
	}
	dockerClient.NegotiateAPIVersion(context.TODO())

	for index := 0; index < numberVolumes; index++ {
		volumeName := fmt.Sprintf("%s-%d", nodePrefix, index)

		if err := dockerClient.VolumeRemove(context.TODO(), volumeName, false); err != nil {
			return fmt.Errorf("removing volume for node: %w", err)
		}
	}

	return nil
}

func testgroup(t *testing.T, groupname string) error {
	t.Logf("testing group %s", groupname)

	version, err := semver.Parse(defaultKubernetesVersion)
	if err != nil {
		return err
	}

	u := uuid.New()

	node := v1alpha3.Node{}
	if err := createNodeVolumes(3, u.String(), &node); err != nil {
		return err
	}
	defer func() {
		if err := cleanupNodeVolumes(3, u.String(), &node); err != nil {
			t.Logf("error: %s", err)
		}
	}()

	cluster, err := kind.NewCluster(version, create.WithV1Alpha3(&v1alpha3.Cluster{
		Nodes: []v1alpha3.Node{ node, },
	}))
	if err != nil {
		return err
	}
	defer cluster.Cleanup()

	if err := temp.DeployController(cluster, "kind"); err != nil {
		return err
	}

	addons, err := addons(addonTestingGroups[groupname]...)
	if err != nil {
		return err
	}

	if err := removeLocalPathAsDefaultStorage(cluster, addons); err != nil {
		return err
	}

	ph, err := test.NewBasicTestHarness(t, cluster, addons...)
	if err != nil {
		return err
	}
	defer ph.Cleanup()

	ph.Validate()
	ph.Deploy()

	return nil
}

func addons(names ...string) ([]v1beta1.AddonInterface, error) {
	var testAddons []v1beta1.AddonInterface

	repo, err := local.NewRepository("base", "../addons")
	if err != nil {
		return testAddons, err
	}
	addons, err := repo.ListAddons()
	if err != nil {
		return testAddons, err
	}

	for _, addon := range addons {
		for _, name := range names {
			overrides(addon[0])
			if addon[0].GetName() == name {
				testAddons = append(testAddons, addon[0])
			}
		}
	}

	if len(testAddons) != len(names) {
		return testAddons, fmt.Errorf("got %d addons, expected %d", len(testAddons), len(names))
	}

	return testAddons, nil
}

func findUnhandled() ([]v1beta1.AddonInterface, error) {
	var unhandled []v1beta1.AddonInterface
	repo, err := local.NewRepository("base", "../addons")
	if err != nil {
		return unhandled, err
	}
	addons, err := repo.ListAddons()
	if err != nil {
		return unhandled, err
	}

	for _, revisions := range addons {
		addon := revisions[0]
		found := false
		for _, v := range addonTestingGroups {
			for _, name := range v {
				if name == addon.GetName() {
					found = true
				}
			}
		}
		if !found {
			unhandled = append(unhandled, addon)
		}
	}

	return unhandled, nil
}

func removeLocalPathAsDefaultStorage(cluster test.Cluster, addons []v1beta1.AddonInterface) error {
	for _, addon := range addons {
		if addon.GetName() == "localvolumeprovisioner" {
			if err := kubectl("--kubeconfig", cluster.ConfigPath(), "patch", "storageclass", "local-path", "-p", patchStorageClass); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func kubectl(args ...string) error {
	cmd := exec.Command("kubectl", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// -----------------------------------------------------------------------------
// Private - CI Values Overrides
// -----------------------------------------------------------------------------

// TODO: a temporary place to put configuration overrides for addons
// See: https://jira.mesosphere.com/browse/DCOS-62137
func overrides(addon v1beta1.AddonInterface) {
	if v, ok := addonOverrides[addon.GetName()]; ok {
		addon.GetAddonSpec().ChartReference.Values = &v
	}
}

var addonOverrides = map[string]string{
	"dispatch": `
---
argo-cd:
  prometheus:
    enabled: true
    release: prometheus-kubeaddons

prometheus:
  enabled: true
  release: prometheus-kubeaddons

minio:
  persistence:
    size: 1Gi
`,
	"metallb": `
---
configInline:
  address-pools:
  - name: default
    protocol: layer2
    addresses:
    - "172.17.1.200-172.17.1.250"
`,
	"istio": `
---
      kiali:
       enabled: true
       contextPath: /ops/portal/kiali
       ingress:
         enabled: true
         kubernetes.io/ingress.class: traefik
         hosts:
           - ""
       dashboard:
         auth:
           strategy: anonymous
       prometheusAddr: http://prometheus-kubeaddons-prom-prometheus.kubeaddons:9090

      tracing:
        enabled: true
        contextPath: /ops/portal/jaeger
        ingress:
          enabled: true
          kubernetes.io/ingress.class: traefik
          hosts:
            - ""

      grafana:
        enabled: true

      prometheus:
        serviceName: prometheus-kubeaddons-prom-prometheus.kubeaddons

      istiocoredns:
        enabled: true

      security:
       selfSigned: true
       caCert: /etc/cacerts/tls.crt
       caKey: /etc/cacerts/tls.key
       rootCert: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
       certChain: /etc/cacerts/tls.crt
       enableNamespacesByDefault: false

      global:
       podDNSSearchNamespaces:
       - global
       - "{{ valueOrDefault .DeploymentMeta.Namespace \"default\" }}.global"

       mtls:
        enabled: true

       multiCluster:
        enabled: true

       controlPlaneSecurityEnabled: true
`,
}
