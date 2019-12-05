package test

import (
	"fmt"
	"testing"

	"github.com/blang/semver"

	"github.com/mesosphere/kubeaddons/api/v1beta1"
	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

const defaultKubernetesVersion = "1.15.6"

var environmentConciousFilteredAddons = []string{
	"awsebscsiprovisioner",
	"awsebsprovisioner",
	"azuredisk-csi-driver",
	"azurediskprovisioner",
	"konvoyconfig",
	"localvolumeprovisioner",
	"metallb",
	"nvidia",
}

// TODO - only doing a couple of addons for the moment, this will be expanded upon in later iterations
// after we've worked out some of the issues with the testing environment and addon requirements.
var temporarilyFilteredAddons = []string{
	"cert-manager",
	"defaultstorageclass-protection",
	"dex-k8s-authenticator",
	"dex",
	"dispatch",
	"external-dns",
	"flagger",
	"gatekeeper",
	"istio",
	"kibana",
	"kommander",
	"kube-oidc-proxy",
	"opsportal",
	"prometheusadapter",
	"prometheus",
	"reloader",
	"traefik-forward-auth",
	"traefik",
	"velero",
	"kudo",
}

// TestAddons tests deployment of all addons in this repository
func TestAddons(t *testing.T) {
	t.Log("testing filtered addon deployment")
	cluster, err := kind.NewCluster(semver.MustParse(defaultKubernetesVersion))
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Cleanup()

	if err := temp.DeployController(cluster); err != nil {
		t.Fatal(err)
	}

	addons, err := temp.Addons("../addons/")
	if err != nil {
		t.Fatal(err)
	}

	var testAddons []v1beta1.AddonInterface
	for _, v := range addons {
		isFiltered := false
		for _, filtered := range append(temporarilyFilteredAddons, environmentConciousFilteredAddons...) {
			if v[0].GetName() == filtered {
				isFiltered = true
			}
		}
		if !isFiltered {
			// TODO - for right now, we're only testing the latest revision.
			// We're waiting on additional features from the test harness to
			// expand this, see https://jira.mesosphere.com/browse/DCOS-61266
			testAddons = append(testAddons, v[0])
		}
	}

	th, err := test.NewBasicTestHarness(t, cluster, testAddons...)
	if err != nil {
		t.Fatal(err)
	}
	defer th.Cleanup()

	th.Validate()
	th.Deploy()
}

func TestElasticSearchDeploy(t *testing.T) {
	t.Log("testing elasticsearch deployment")
	cluster, err := kind.NewCluster(semver.MustParse(defaultKubernetesVersion))
	if err != nil {
		t.Fatal(err)
	}
	defer cluster.Cleanup()

	if err := temp.DeployController(cluster); err != nil {
		t.Fatal(err)
	}

	addons, err := addons("elasticsearch", "elasticsearchexporter", "kibana")
	if err != nil {
		t.Fatal(err)
	}

	ph, err := test.NewBasicTestHarness(t, cluster, addons...)
	if err != nil {
		t.Fatal(err)
	}
	defer ph.Cleanup()

	ph.Validate()
	ph.Deploy()

}

func TestPrometheusDeploy(t *testing.T) {
	t.Log("testing prometheus deployment")
	promCluster, err := kind.NewCluster(semver.MustParse(defaultKubernetesVersion))
	if err != nil {
		t.Fatal(err)
	}
	defer promCluster.Cleanup()

	if err := temp.DeployController(promCluster); err != nil {
		t.Fatal(err)
	}

	addons, err := addons("prometheus", "prometheus-adapter", "opsportal")
	if err != nil {
		t.Fatal(err)
	}

	ph, err := test.NewBasicTestHarness(t, promCluster, addons...)
	if err != nil {
		t.Fatal(err)
	}
	defer ph.Cleanup()

	ph.Validate()
	ph.Deploy()
}

// -----------------------------------------------------------------------------
// Private Functions
// -----------------------------------------------------------------------------

func addons(names ...string) ([]v1beta1.AddonInterface, error) {
	var testAddons []v1beta1.AddonInterface

	addons, err := temp.Addons("../addons/")
	if err != nil {
		return testAddons, err
	}

	for _, addon := range addons {
		for _, name := range names {
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
