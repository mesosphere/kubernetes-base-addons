package test

import (
	"testing"

	"github.com/blang/semver"

	"github.com/mesosphere/kubeaddons/api/v1beta1"
	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

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
	cluster, err := kind.NewCluster(semver.MustParse("1.15.6"))
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

func TestPrometheusDeploy(t *testing.T){
	t.Log("testing prometheus deployment")
	// test prometheus
	promCluster, err := kind.NewCluster(semver.MustParse("1.15.6"))
	if err != nil {
		t.Fatal(err)
	}
	defer promCluster.Cleanup()

	if err := temp.DeployController(promCluster); err != nil {
		t.Fatal(err)
	}

	addons, err := temp.Addons("../addons/")
	if err != nil {
		t.Fatal(err)
	}

	var testAddons []v1beta1.AddonInterface
	for _, addon := range addons {
		switch addon[0].GetName() {
		case "prometheus", "prometheusadapter", "opsportal":
			testAddons = append(testAddons, addon[0])
		default:
			continue
		}
	}

	ph, err := test.NewBasicTestHarness(t, promCluster, testAddons...)
	if err != nil {
		t.Fatal(err)
	}
	defer ph.Cleanup()

	ph.Validate()
	ph.Deploy()
}
