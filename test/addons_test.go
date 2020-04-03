package test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/blang/semver"
	volumetypes "github.com/docker/docker/api/types/volume"
	docker "github.com/docker/docker/client"
	"github.com/google/uuid"
	testcluster "github.com/mesosphere/ksphere-testing-framework/pkg/cluster"
	"github.com/mesosphere/ksphere-testing-framework/pkg/cluster/kind"
	"github.com/mesosphere/ksphere-testing-framework/pkg/cluster/konvoy"
	"github.com/mesosphere/ksphere-testing-framework/pkg/experimental"
	testharness "github.com/mesosphere/ksphere-testing-framework/pkg/harness"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/catalog"
	"github.com/mesosphere/kubeaddons/pkg/repositories"
	"github.com/mesosphere/kubeaddons/pkg/repositories/git"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
	addontesters "github.com/mesosphere/kubeaddons/test/utils"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha3"
	"sigs.k8s.io/kind/pkg/cluster"
)

const (
	controllerBundle         = "https://mesosphere.github.io/kubeaddons/bundle.yaml"
	defaultKubernetesVersion = "1.15.6"
	patchStorageClass        = `{"metadata": {"annotations":{"storageclass.kubernetes.io/is-default-class":"false"}}}`

	comRepoURL    = "https://github.com/mesosphere/kubeaddons-community"
	comRepoRef    = "master"
	comRepoRemote = "origin"
)

var (
	cat       catalog.Catalog
	localRepo repositories.Repository
	comRepo   repositories.Repository
	groups    map[string][]v1beta1.AddonInterface
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
	if err := testgroup(t, "general"); err != nil {
		t.Fatal(err)
	}
}

func TestBackupsGroup(t *testing.T) {
	if err := testgroup(t, "backups"); err != nil {
		t.Fatal(err)
	}
}

func TestSsoGroup(t *testing.T) {
	if err := testgroup(t, "sso"); err != nil {
		t.Fatal(err)
	}
}

func TestElasticsearchGroup(t *testing.T) {
	if err := testgroup(t, "elasticsearch", elasticsearchChecker, kibanaChecker); err != nil {
		t.Fatal(err)
	}
}

func TestPrometheusGroup(t *testing.T) {
	if err := testgroup(t, "prometheus"); err != nil {
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

func TestAwsGroup(t *testing.T) {
	if err := testgroup(t, "aws"); err != nil {
		t.Fatal(err)
	}
}

func TestAzureGroup(t *testing.T) {
	if err := testgroup(t, "azure"); err != nil {
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

func testgroup(t *testing.T, groupname string, jobs ...clusterTestJob) error {
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

	var tcluster testcluster.Cluster
	if groupname == "aws" || groupname == "azure" || groupname == "gcp" {
		path, _ := os.Getwd()
		tcluster, err = konvoy.NewCluster(fmt.Sprintf("%s/konvoy", path), groupname)
	} else {
		tcluster, err = kind.NewCluster(version, cluster.CreateWithV1Alpha3Config(&v1alpha3.Cluster{Nodes: []v1alpha3.Node{node}}))
	}
	if err != nil {
		// try to clean up in case cluster was created and reference available
		if tcluster != nil {
			_ = tcluster.Cleanup()
		}
		return err
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
		overrides(addon)
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
func overrides(addon v1beta1.AddonInterface) {
	if v, ok := addonOverrides[addon.GetName()]; ok {
		addon.GetAddonSpec().ChartReference.Values = &v
	}
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
	// Copied from addons/elasticsearch/6.8.x/elasticsearch-5.yaml and edited to lower resource
	// limits so it can deploy on a kind cluster.
	// Added and edited lines are marked with "# override".
	"elasticsearch": `
---
client:
  heapSize: 256m  # override
  resources:
    limits:
      cpu: 1000m  # override
      memory: 512Mi  # override
    requests:
      cpu: 500m  # override
      memory: 256Mi  # override
master:
  updateStrategy:
    type: RollingUpdate
  heapSize: 256m  # override
  resources:
    # need more cpu upon initialization, therefore burstable class
    limits:
      cpu: 1000m  # override
      memory: 512Mi  # override
    requests:
      cpu: 100m  # override
      memory: 256Mi  # override
data:
  updateStrategy:
    type: RollingUpdate
  hooks:
    drain:
      enabled: false
    # Because the drain is set to false, we can take advantage here and create resources we need
    postStart: |-
      #!/bin/bash
      set -o errexit
      set -o nounset
      set -o pipefail
      
      # Wait until client nodes accept requests.
      # This prevents data pods from getting into a long CrashLoopBackoff.
      echo "Waiting for client node..."
      # Use a for loop to retry on connection failures.
      # Elasticsearch image's curl doesn't have --retry-connrefused.
      for i in {1..10}; do curl --fail '{{ template "elasticsearch.client.fullname" . }}:9200/_cluster/health' && s=0 && break || s=$? && sleep 5; done
      if [ $s -ne 0 ]; then echo "client node not available. Exiting"; exit $s; fi
      
      # Creating the index template: 'kubernetes_cluster'
      # Reduces the number of fields produced due to the indexing of the audit logs
      # if template doesnt return 200, try to create it
      if ! curl -I -XGET '{{ template "elasticsearch.client.fullname" . }}:9200/_template/kubernetes_cluster' | grep "200" > /dev/null
      then
          echo "Creating the index template: 'kubernetes_cluster'"
          curl -H 'Content-Type: application/json' -XPUT '{{ template "elasticsearch.client.fullname" . }}:9200/_template/kubernetes_cluster' -d '{"index_patterns":["kubernetes_cluster*"],"mappings":{"flb_type":{"properties":{"@ts":{"type":"date"},"requestObject":{"dynamic":false,"properties":{"apiVersion":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"kind":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"metadata":{"properties":{"creationTimestamp":{"type":"date"},"labels":{"properties":{"app":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"app_kubernetes_io/instance":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"app_kubernetes_io/managed-by":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"app_kubernetes_io/name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"app_kubernetes_io/version":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"helm_sh/chart":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"namespace":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"ownerReferences":{"properties":{"apiVersion":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"kind":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"uid":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"uid":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"status":{"properties":{"allowed":{"type":"boolean"},"conditions":{"properties":{"lastHeartbeatTime":{"type":"date"},"lastTransitionTime":{"type":"date"},"lastUpdateTime":{"type":"date"},"message":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"reason":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"status":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"type":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"containerStatuses":{"properties":{"containerID":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"image":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"imageID":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"lastState":{"properties":{"terminated":{"properties":{"containerID":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"exitCode":{"type":"long"},"finishedAt":{"type":"date"},"reason":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"startedAt":{"type":"date"}}}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"ready":{"type":"boolean"},"restartCount":{"type":"long"},"started":{"type":"boolean"},"state":{"properties":{"running":{"properties":{"startedAt":{"type":"date"}}},"waiting":{"properties":{"reason":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}}}},"hostIP":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"loadBalancer":{"properties":{"ingress":{"properties":{"hostname":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}},"podIP":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"spec":{"properties":{"user":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"template":{"properties":{"spec":{"properties":{"containers":{"properties":{"args":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"command":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"env":{"properties":{"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"value":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"valueFrom":{"properties":{"fieldRef":{"properties":{"apiVersion":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"fieldPath":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}}}},"image":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"imagePullPolicy":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"livenessProbe":{"properties":{"failureThreshold":{"type":"long"},"httpGet":{"properties":{"path":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"port":{"type":"long"},"scheme":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"initialDelaySeconds":{"type":"long"},"periodSeconds":{"type":"long"},"successThreshold":{"type":"long"},"timeoutSeconds":{"type":"long"}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"ports":{"properties":{"containerPort":{"type":"long"},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"protocol":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"readinessProbe":{"properties":{"failureThreshold":{"type":"long"},"httpGet":{"properties":{"path":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"port":{"type":"long"},"scheme":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"initialDelaySeconds":{"type":"long"},"periodSeconds":{"type":"long"},"successThreshold":{"type":"long"},"timeoutSeconds":{"type":"long"}}},"resources":{"properties":{"limits":{"properties":{"cpu":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"memory":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"requests":{"properties":{"cpu":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"memory":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}},"terminationMessagePath":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"terminationMessagePolicy":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"volumeMounts":{"properties":{"mountPath":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"readOnly":{"type":"boolean"}}}}},"serviceAccount":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"serviceAccountName":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}}}},"webhooks":{"properties":{"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"namespaceSelector":{"properties":{"matchExpressions":{"properties":{"key":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"operator":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"values":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}},"objectSelector":{"properties":{"matchExpressions":{"properties":{"key":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"operator":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"values":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}}}}}},"responseObject":{"dynamic":false,"properties":{"apiVersion":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"kind":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"metadata":{"properties":{"creationTimestamp":{"type":"date"},"labels":{"properties":{"app":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"app_kubernetes_io/instance":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"app_kubernetes_io/managed-by":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"app_kubernetes_io/name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"app_kubernetes_io/version":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"chart":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"helm_sh/chart":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"namespace":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"ownerReferences":{"properties":{"apiVersion":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"kind":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"uid":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"uid":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"status":{"properties":{"allowed":{"type":"boolean"},"conditions":{"properties":{"lastTransitionTime":{"type":"date"},"lastUpdateTime":{"type":"date"},"message":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"reason":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"status":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"type":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"reason":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"spec":{"properties":{"template":{"properties":{"spec":{"properties":{"containers":{"properties":{"command":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"image":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"imagePullPolicy":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"ports":{"properties":{"containerPort":{"type":"long"},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"protocol":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"volumeMounts":{"properties":{"mountPath":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"readOnly":{"type":"boolean"}}}}},"serviceAccount":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"serviceAccountName":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}},"user":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}},"webhooks":{"properties":{"name":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"namespaceSelector":{"properties":{"matchExpressions":{"properties":{"key":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"operator":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"values":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}}}}}},"responseStatus":{"dynamic":false,"properties":{"code":{"type":"long"},"message":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"metadata":{"type":"object"},"reason":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}},"status":{"type":"text","fields":{"keyword":{"type":"keyword","ignore_above":256}}}}}}}}}'
      fi
  persistence:  # override
    size: 4Gi  # override
  heapSize: 1024m  # override
  resources:
    # need more cpu upon initialization, therefore burstable class
    limits:
      cpu: 1000m  # override
      memory: 1536Mi  # override
    requests:
      cpu: 100m  # override
      memory: 1024Mi  # override
`,
}
