package test

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/blang/semver"
	volumetypes "github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
	"github.com/google/uuid"
	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	"github.com/mesosphere/ksphere-testing-framework/pkg/cluster/kind"
	"github.com/mesosphere/ksphere-testing-framework/pkg/cluster/konvoy"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta2"
	"github.com/mesosphere/kubeaddons/pkg/catalog"
	"github.com/mesosphere/kubeaddons/pkg/constants"
	"github.com/mesosphere/kubeaddons/pkg/repositories"
	"github.com/mesosphere/kubeaddons/pkg/repositories/git"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
	testutils "github.com/mesosphere/kubeaddons/test/utils"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/helm/pkg/chartutil"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
)

const (
	controllerBundle        = "https://mesosphere.github.io/kubeaddons/bundle.yaml"
	defaultKindestNodeImage = "kindest/node:v1.18.8"
	patchStorageClass       = `{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}`

	comRepoURL    = "https://github.com/mesosphere/kubeaddons-community"
	comRepoRemote = "origin"
	comRepoRef    = "master"

	defaultKBARepoRef = "master"

	allAWSGroupName = "allAWS"

	tempDir = "/tmp/kubernetes-base-addons"

	kubeaddonsControllerNamespace = "kubeaddons"
	kubeaddonsControllerPodPrefix = "kubeaddons-controller-manager-"
)

var (
	cat        catalog.Catalog
	localRepo  repositories.Repository
	comRepo    repositories.Repository
	kbaRepoRef string
	groups     map[string][]v1beta2.AddonInterface
)

type clusterTestJob func(*testing.T, testcluster.Cluster) testharness.Job

var kbaBranchFlag = flag.String("kba-branch", "", "")

func TestMain(m *testing.M) {
	flag.Parse()

	if *kbaBranchFlag != "" {
		kbaRepoRef = *kbaBranchFlag
	} else {
		kbaRepoRef = defaultKBARepoRef
	}

	var err error

	fmt.Println("initializing local repository for test...")
	localRepo, err = local.NewRepository("local", "../addons")
	if err != nil {
		panic(err)
	}

	fmt.Printf("initializing remote repository %s for test...\n", comRepoURL)
	comRepo, err = git.NewRemoteRepository(comRepoURL, comRepoRef, comRepoRemote)
	if err != nil {
		panic(err)
	}

	fmt.Println("initializing catalog with repositories...")
	cat, err = catalog.NewCatalog(localRepo, comRepo)
	if err != nil {
		panic(err)
	}

	fmt.Println("finding addon test groups...")
	groupsMap, err := getGroupsMapFromFile("groups.yaml")
	if err != nil {
		panic(err)
	}
	appendDynamicToGroupsMap(groupsMap)
	groups, err = testutils.AddonsForGroups(groupsMap, cat)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func getGroupsMapFromFile(f string) (testutils.Groups, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}

	g := make(testutils.Groups)
	if err := yaml.Unmarshal(b, &g); err != nil {
		return nil, err
	}
	return g, nil
}

// appends all AWS related addons to the groupsMap as allAWSGroupName
func appendDynamicToGroupsMap(groupsMap testutils.Groups) {
	addonRevisionsList, err := cat.ListAddons(func(addon v1beta2.AddonInterface) bool {
		// https://github.com/mesosphere/konvoy/blob/94899699aa49ce8344a9d000300d9fa37ebbbf48/pkg/addons/addons.go#L97-L99
		if len(addon.GetAddonSpec().CloudProvider) == 0 {
			return true
		}
		for _, cloudProvider := range addon.GetAddonSpec().CloudProvider {
			if cloudProvider.Name == "aws" && cloudProvider.Enabled {
				return true
			}
		}
		return false
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nAdding the following addons to Dynamic Group, %v:\n", allAWSGroupName)
	for _, addonRevisons := range addonRevisionsList {
		addonName := addonRevisons[0].GetName()
		fmt.Println(addonName)
		groupsMap[allAWSGroupName] = append(groupsMap[allAWSGroupName], testutils.AddonName(addonName))
	}
	fmt.Println("")
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

func TestAmbassadorGroup(t *testing.T) {
	if err := testgroup(t, "ambassador", defaultKindestNodeImage, ambassadorChecker); err != nil {
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

func checkIfUpgradeIsNeeded(t *testing.T, groupname string) (bool, []v1beta2.AddonInterface, error) {
	t.Logf("determining which addons in group %s need to be upgrade tested", groupname)

	var doUpgrade bool
	addons := groups[groupname]
	for _, addon := range addons {
		if err := overrides(addon); err != nil {
			return false, nil, err
		}
	}

	addonDeploymentsArray := make([]v1beta2.AddonInterface, 0)
	for _, newAddon := range addons {
		t.Logf("verifying whether upgrade testing is needed for addon %s", newAddon.GetName())
		oldAddon, err := testutils.GetLatestAddonRevisionFromLocalRepoBranch("../", comRepoRemote, kbaRepoRef, newAddon.GetName())
		if err != nil {
			if strings.Contains(err.Error(), "directory not found") {
				t.Logf("no need to upgrade test %s, it appears to be a new addon (no previous revisions found in branch %s)", newAddon.GetName(), kbaRepoRef)
				addonDeploymentsArray = append(addonDeploymentsArray, newAddon)
				continue
			}
			return false, nil, err
		}
		if oldAddon == nil {
			t.Logf("no need to upgrade test %s, it appears to be a new addon (no previous revisions found in branch %s)", newAddon.GetName(), kbaRepoRef)
			addonDeploymentsArray = append(addonDeploymentsArray, newAddon)
			continue // new addon, upgrade test not needed
		}

		// Apply overrides to oldAddon to ensure it is deployed with the necessary value overrides
		if err := overrides(oldAddon); err != nil {
			return false, nil, err
		}

		t.Logf("determining old and new versions for upgrade testing addon %s", newAddon.GetName())
		oldRev := oldAddon.GetAnnotations()[constants.AddonRevisionAnnotation]
		oldVersion, err := semver.Parse(strings.TrimPrefix(oldRev, "v"))
		if err != nil {
			return false, nil, err
		}
		newRev := newAddon.GetAnnotations()[constants.AddonRevisionAnnotation]
		newVersion, err := semver.Parse(strings.TrimPrefix(newRev, "v"))
		if err != nil {
			return false, nil, err
		}

		if newVersion.EQ(oldVersion) {
			t.Logf("skipping upgrade test for addon %s, it has not been updated", newAddon.GetName())
			addonDeploymentsArray = append(addonDeploymentsArray, oldAddon)
			continue
		} else if oldVersion.GT(newVersion) {
			return false, nil, fmt.Errorf("revisions for addon %s are broken, previous revision %s is newer than current %s", newAddon.GetName(), oldVersion, newVersion)
		}

		t.Logf("found old version of addon %s %s (revision %s) and new version %s (revision %s)", newAddon.GetName(), oldRev, oldVersion, newVersion, newRev)
		// for upgraded addons, add the oldAddon (running previous version) to deployments
		addonDeploymentsArray = append(addonDeploymentsArray, oldAddon)
		doUpgrade = true
	}

	return doUpgrade, addonDeploymentsArray, nil
}

func createNodeVolumes(numberVolumes int, nodePrefix string, node *v1alpha4.Node) error {
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

		node.ExtraMounts = append(node.ExtraMounts, v1alpha4.Mount{
			ContainerPath: fmt.Sprintf("/mnt/disks/%s", volumeName),
			HostPath:      volume.Mountpoint,
		})
	}

	return nil
}

func cleanupNodeVolumes(numberVolumes int, nodePrefix string, node *v1alpha4.Node) error {
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

	node := v1alpha4.Node{}
	if err := createNodeVolumes(3, u.String(), &node); err != nil {
		return err
	}
	defer func() {
		if err := cleanupNodeVolumes(3, u.String(), &node); err != nil {
			t.Logf("error: %s", err)
		}
	}()

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}
	dir, err := ioutil.TempDir(tempDir, groupname+"-")
	if err != nil {
		return err
	}

	t.Logf("setting up cluster for test group %s", groupname)
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

	kubeConfig, err := tcluster.ConfigYAML()
	if err != nil {
		return err
	}

	kubeConfigPath := filepath.Join(dir, "kubeconfig")
	if err := ioutil.WriteFile(kubeConfigPath, kubeConfig, 0644); err != nil {
		return err
	}

	if err := kubectl("--kubeconfig", kubeConfigPath, "apply", "-f", controllerBundle); err != nil {
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
	go testutils.LoggingHook(t, tcluster, wg, stop)

	addonDeployment, err := testutils.DeployAddons(t, tcluster, addons...)
	if err != nil {
		return err
	}

	addonCleanup, err := testutils.CleanupAddons(t, tcluster, addons...)
	if err != nil {
		return err
	}

	addonDefaults, err := testutils.WaitForAddons(t, tcluster, addons...)
	if err != nil {
		return err
	}

	t.Logf("determining which addons in group %s need to be upgrade tested", groupname)
	addonUpgrades := testharness.Loadables{}
	for _, newAddon := range addons {
		t.Logf("verifying whether upgrade testing is needed for addon %s", newAddon.GetName())
		oldAddon, err := testutils.GetLatestAddonRevisionFromLocalRepoBranch("../", comRepoRemote, kbaRepoRef, newAddon.GetName())
		if err != nil {
			if strings.Contains(err.Error(), "directory not found") {
				t.Logf("no need to upgrade test %s, it appears to be a new addon (no previous revisions found in branch %s)", newAddon.GetName(), kbaRepoRef)
				continue
			}
			return err
		}
		if oldAddon == nil {
			t.Logf("no need to upgrade test %s, it appears to be a new addon (no previous revisions found in branch %s)", newAddon.GetName(), kbaRepoRef)
			continue // new addon, upgrade test not needed
		}

		t.Logf("determining old and new versions for upgrade testing addon %s", newAddon.GetName())
		oldRev := oldAddon.GetAnnotations()[constants.AddonRevisionAnnotation]
		oldVersion, err := semver.Parse(strings.TrimPrefix(oldRev, "v"))
		if err != nil {
			return err
		}
		newRev := newAddon.GetAnnotations()[constants.AddonRevisionAnnotation]
		newVersion, err := semver.Parse(strings.TrimPrefix(newRev, "v"))
		if err != nil {
			return err
		}

		t.Logf("found old version of addon %s %s (revision %s) and new version %s (revision %s)", newAddon.GetName(), oldRev, oldVersion, newVersion, newRev)
		if oldVersion.GT(newVersion) {
			return fmt.Errorf("revisions for addon %s are broken, previous revision %s is newer than current %s", newAddon.GetName(), oldVersion, newVersion)
		}
		if !newVersion.GT(oldVersion) {
			t.Logf("skipping upgrade test for addon %s, it has not been updated", newAddon.GetName())
			continue
		}

		t.Logf("INFO: addon %s was modified and will be upgrade tested", newAddon.GetName())
		addonUpgrade, err := testutils.UpgradeAddon(t, tcluster, oldAddon, newAddon)
		if err != nil {
			return err
		}

		addonUpgrades = append(addonUpgrades, addonUpgrade)
	}

	th := testharness.NewSimpleTestHarness(t)
	th.Load(
		testutils.ValidateAddons(addons...),
		addonDeployment,
		addonDefaults,
	)
	th.Load(addonUpgrades...)

	for _, job := range jobs {
		th.Load(testharness.Loadable{
			Plan: testharness.DefaultPlan,
			Jobs: testharness.Jobs{job(t, tcluster)},
		})
	}

	th.Load(addonCleanup)

	// Collect kubeaddons controller logs during cleanup.
	th.Load(testharness.Loadable{
		Plan: testharness.CleanupPlan,
		Jobs: testharness.Jobs{func(t *testing.T) error {
			logFilePath := filepath.Join(dir, "kubeaddons-controller-log.txt")
			t.Logf("INFO: writing kubeaddons controller logs to %s", logFilePath)

			logFile, err := os.Create(logFilePath)
			if err != nil {
				return err
			}
			defer logFile.Close()

			logs, err := logsFromPodWithPrefix(tcluster, kubeaddonsControllerNamespace, kubeaddonsControllerPodPrefix)
			if err != nil {
				return err
			}
			defer logs.Close()

			_, err = io.Copy(logFile, logs)
			return err
		}},
	})

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

func newCluster(groupname string, version string, node v1alpha4.Node, t *testing.T) (testcluster.Cluster, error) {
	if groupname == "aws" || groupname == "azure" || groupname == "gcp" || groupname == allAWSGroupName {
		provisioner := groupname
		if groupname == allAWSGroupName {
			provisioner = "aws"
		}
		path, _ := os.Getwd()
		return konvoy.NewCluster(fmt.Sprintf("%s/konvoy", path), provisioner)
	}

	path, ok := os.LookupEnv("KBA_KUBECONFIG")
	if !ok {
		t.Logf("No Kubeconfig specified in KBA_KUBECONFIG. Creating Kind cluster")
		return kind.NewCluster(version, cluster.CreateWithV1Alpha4Config(&v1alpha4.Cluster{Nodes: []v1alpha4.Node{node}}))
	}

	config, err := clientcmd.LoadFromFile(path)
	if err != nil || len(config.Contexts) == 0 {
		t.Logf("%s is not a valid kubeconfig. Creating Kind cluster", path)
		return kind.NewCluster(version, cluster.CreateWithV1Alpha4Config(&v1alpha4.Cluster{Nodes: []v1alpha4.Node{node}}), cluster.CreateWithKubeconfigPath(path))
	}

	t.Log("Using KBA_KUBECONFIG at", path)
	// load the file from kubeconfig
	kubeConfig, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return testcluster.NewClusterFromKubeConfig("kind", kubeConfig)
}
