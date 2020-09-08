package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/mesosphere/ksphere-testing-framework/pkg/cluster/kind"
	"sigs.k8s.io/kind/pkg/cluster"

	volumetypes "github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
	"github.com/google/uuid"
	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	"github.com/mesosphere/ksphere-testing-framework/pkg/cluster/konvoy"
	"github.com/mesosphere/ksphere-testing-framework/pkg/experimental"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta2"
	"github.com/mesosphere/kubeaddons/pkg/catalog"
	"github.com/mesosphere/kubeaddons/pkg/repositories"
	"github.com/mesosphere/kubeaddons/pkg/repositories/git"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
	addontesters "github.com/mesosphere/kubeaddons/test/utils"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/chartutil"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha3"
)

const (
	controllerBundle = "https://mesosphere.github.io/kubeaddons/bundle.yaml"
	// just take the default from ksphere-testing-framework with ""
	defaultKindestNodeImage = "kindest/node:v1.18.6@sha256:b9f76dd2d7479edcfad9b4f636077c606e1033a2faf54a8e1dee6509794ce87d"
	patchStorageClass       = `{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}`

	comRepoURL    = "https://github.com/mesosphere/kubeaddons-community"
	comRepoRef    = "master"
	comRepoRemote = "origin"
)

var (
	cat       catalog.Catalog
	localRepo repositories.Repository
	comRepo   repositories.Repository
	groups    map[string][]v1beta2.AddonInterface
)

type clusterTestJob func(*testing.T, testcluster.Cluster) testharness.Job

func init() {
	var err error

	localRepo, err = local.NewRepository("local", "../addons")
	if err != nil {
		panic(err)
	}

	comRepo, err = git.NewRemoteRepository(comRepoURL, comRepoRef, comRepoRemote)
	if err != nil {
		panic(err)
	}

	cat, err = catalog.NewCatalog(localRepo, comRepo)
	if err != nil {
		panic(err)
	}

	groups, err = experimental.AddonsForGroupsFile("groups.yaml", cat)
	if err != nil {
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
	if err := testgroup(t, "general", defaultKindestNodeImage); err != nil {
		t.Fatal(err)
	}
}

func TestBackupsGroup(t *testing.T) {
	if err := testgroup(t, "backups", defaultKindestNodeImage); err != nil {
		t.Fatal(err)
	}
}

func TestSsoGroup(t *testing.T) {
	if err := testgroup(t, "sso", defaultKindestNodeImage); err != nil {
		t.Fatal(err)
	}
}

func TestElasticsearchGroup(t *testing.T) {
	if err := testgroup(t, "elasticsearch", defaultKindestNodeImage, elasticsearchChecker, kibanaChecker); err != nil {
		t.Fatal(err)
	}
}

func TestPrometheusGroup(t *testing.T) {
	if err := testgroup(t, "prometheus", defaultKindestNodeImage, promChecker, alertmanagerChecker, grafanaChecker); err != nil {
		t.Fatal(err)
	}
}

func TestIstioGroup(t *testing.T) {
	if err := testgroup(t, "istio", "kindest/node:v1.16.9"); err != nil {
		t.Fatal(err)
	}
}

func TestLocalVolumeProvisionerGroup(t *testing.T) {
	if err := testgroup(t, "localvolumeprovisioner", defaultKindestNodeImage); err != nil {
		t.Fatal(err)
	}
}

func TestAwsGroup(t *testing.T) {
	if err := testgroup(t, "aws", defaultKindestNodeImage); err != nil {
		t.Fatal(err)
	}
}

func TestAzureGroup(t *testing.T) {
	if err := testgroup(t, "azure", defaultKindestNodeImage); err != nil {
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

		node.ExtraMounts = append(node.ExtraMounts, v1alpha3.Mount{
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

func testgroup(t *testing.T, groupname string, version string, jobs ...clusterTestJob) error {
	var err error
	t.Logf("testing group %s", groupname)

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

	tcluster, err := newCluster(groupname, version, node, t)
	if err != nil {
		// try to clean up in case cluster was created and reference available
		if tcluster != nil {
			_ = tcluster.Cleanup()
		}
		return err
	}
	if tcluster == nil {
		return fmt.Errorf("tcluster is nil")
	}
	defer tcluster.Cleanup()

	f, err := ioutil.TempFile(os.TempDir(), "konvoy-test-")
	if err != nil {
		return err
	}

	kubeConfig, err := tcluster.ConfigYAML()
	if err != nil {
		return err
	}

	if _, err := f.Write(kubeConfig); err != nil {
		return err
	}

	if err := kubectl("--kubeconfig", f.Name(), "apply", "-f", controllerBundle); err != nil {
		return err
	}

	addons := groups[groupname]
	for _, addon := range addons {
		if err := overrides(addon); err != nil {
			return err
		}
	}

	wg := &sync.WaitGroup{}
	stop := make(chan struct{})
	go experimental.LoggingHook(t, tcluster, wg, stop)

	addonDeployment, err := addontesters.DeployAddons(t, tcluster, addons...)
	if err != nil {
		return err
	}

	addonCleanup, err := addontesters.CleanupAddons(t, tcluster, addons...)
	if err != nil {
		return err
	}

	addonDefaults, err := addontesters.WaitForAddons(t, tcluster, addons...)
	if err != nil {
		return err
	}

	th := testharness.NewSimpleTestHarness(t)
	th.Load(
		addontesters.ValidateAddons(addons...),
		addonDeployment,
		addonDefaults,
		addonCleanup,
	)
	for _, job := range jobs {
		th.Load(testharness.Loadable{
			Plan: testharness.DefaultPlan,
			Jobs: testharness.Jobs{job(t, tcluster)},
		})
	}

	defer th.Cleanup()
	th.Validate()
	th.Deploy()
	th.Default()

	close(stop)
	wg.Wait()

	return nil
}

func findUnhandled() ([]v1beta2.AddonInterface, error) {
	var unhandled []v1beta2.AddonInterface
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
		for _, v := range groups {
			for _, r := range v {
				if r.GetName() == addon.GetName() {
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
func overrides(addon v1beta2.AddonInterface) error {
	overrideValues, ok := addonOverrides[addon.GetName()]
	if !ok {
		return nil
	}

	base := ""
	if addon.GetAddonSpec().ChartReference != nil && addon.GetAddonSpec().ChartReference.Values != nil {
		base = *addon.GetAddonSpec().ChartReference.Values
	}

	values, err := chartutil.ReadValues([]byte(base))
	if err != nil {
		return fmt.Errorf("error decoding values from Addon %s: %v", addon.GetName(), err)
	}

	overrides, err := chartutil.ReadValues([]byte(overrideValues))
	if err != nil {
		return fmt.Errorf("error decoding override values for Addon %s: %v", addon.GetName(), err)
	}

	values.MergeInto(overrides)
	mergedValues, err := values.YAML()
	if err != nil {
		return fmt.Errorf("error merging override values with Addon values for %s: %v", addon.GetName(), err)
	}

	addon.GetAddonSpec().ChartReference.Values = &mergedValues
	return nil
}

var addonOverrides = map[string]string{
	"metallb": `
---
configInline:
  address-pools:
  - name: default
    protocol: layer2
    addresses:
    - "172.17.1.200-172.17.1.250"
`,
	"elasticsearch": `
---
# Reduce resource limits so elasticsearch will deploy on a kind cluster with limited memory.
client:
  heapSize: 256m
  resources:
    limits:
      cpu: 1000m
      memory: 512Mi
    requests:
      cpu: 500m
      memory: 256Mi
master:
  heapSize: 256m
  resources:
    limits:
      cpu: 1000m
      memory: 512Mi
    requests:
      cpu: 100m
      memory: 256Mi
data:
  persistence:
    size: 4Gi
  heapSize: 1024m
  resources:
    limits:
      cpu: 1000m
      memory: 1536Mi
    requests:
      cpu: 100m
      memory: 1024Mi
`,
	"prometheus": `
---
# Remove dependency on persistent volumes and Konvoy's "etcd-certs" secret.
prometheus:
  prometheusSpec:
    secrets: []
    storageSpec: null
kubeEtcd:
  enabled: false
`,
}

func newCluster(groupname string, version string, node v1alpha3.Node, t *testing.T) (testcluster.Cluster, error) {
	if groupname == "aws" || groupname == "azure" || groupname == "gcp" {
		path, _ := os.Getwd()
		return konvoy.NewCluster(fmt.Sprintf("%s/konvoy", path), groupname)
	}

	path, ok := os.LookupEnv("KBA_KUBECONFIG")
	if !ok {
		t.Logf("No Kubeconfig specified in KBA_KUBECONFIG. Creating Kind cluster")
		return kind.NewCluster(version, cluster.CreateWithV1Alpha3Config(&v1alpha3.Cluster{Nodes: []v1alpha3.Node{node}}))
	}

	config, err := clientcmd.LoadFromFile(path)
	if err != nil || len(config.Contexts) == 0 {
		t.Logf("%s is not a valid kubeconfig. Creating Kind cluster", path)
		return kind.NewCluster(version, cluster.CreateWithV1Alpha3Config(&v1alpha3.Cluster{Nodes: []v1alpha3.Node{node}}), cluster.CreateWithKubeconfigPath(path))
	}

	t.Log("Using KBA_KUBECONFIG at", path)
	// load the file from kubeconfig
	kubeConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return testcluster.NewClusterFromKubeConfig("kind", kubeConfig)
}
