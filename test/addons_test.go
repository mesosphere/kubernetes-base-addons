package test

import (
	"fmt"
	"testing"

	"github.com/blang/semver"

	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

const defaultKubernetesVersion = "1.15.6"

var addonTestingGroups = map[string][]string{
	// general - put smaller scope, low resource addons here to be tested in batch
	"general": []string{"dashboard", "external-dns"},

	// elasticsearch - put logging addons which rely on elasticsearch here
	"elasticsearch": []string{"elasticsearch", "elasticsearchexporter", "kibana", "fluentbit"},

	// prometheus - put monitoring addons which rely on prometheus here
	"prometheus": []string{"prometheus", "prometheusadapter", "opsportal"},
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

// -----------------------------------------------------------------------------
// Private Functions
// -----------------------------------------------------------------------------

func testgroup(t *testing.T, groupname string) error {
	t.Logf("testing group %s", groupname)
	cluster, err := kind.NewCluster(semver.MustParse(defaultKubernetesVersion))
	if err != nil {
		return err
	}
	defer cluster.Cleanup()

	if err := temp.DeployController(cluster); err != nil {
		return err
	}

	addons, err := addons(addonTestingGroups[groupname]...)
	if err != nil {
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

var disabled = []string{
	// kudo gets tested in https://github.com/mesosphere/kubeaddons-enterprise, and is likely going to be removed from this repository.
	// See: https://jira.mesosphere.com/browse/DCOS-61842
	"kudo",

	// the following addons need tests added
	// See: https://jira.mesosphere.com/browse/DCOS-61664
	"cert-manager",
	"dex-k8s-authenticator",
	"kube-oidc-proxy",
	"prometheusadapter",
	"velero",
	"dispatch",
	"kommander",
	"traefik",
	"dex",
	"traefik-forward-auth",
	"istio",
	"flagger",
	"gatekeeper",
	"reloader",
	"localvolumeprovisioner",
	"defaultstorageclass-protection",
}

// environmentConciousFilteredAddons are addons which are currently filtered out of tests because we're waiting on features to be able to test them properly.
// See: https://jira.mesosphere.com/browse/DCOS-61664
var environmentConciousFilteredAddons = []string{
	"dex",
	"dex-k8s-authenticator",
	"awsebscsiprovisioner",
	"awsebsprovisioner",
	"azuredisk-csi-driver",
	"azurediskprovisioner",
	"konvoyconfig",
	"localvolumeprovisioner",
	"metallb",
	"nvidia",
}

func findUnhandled() ([]v1beta1.AddonInterface, error) {
	var unhandled []v1beta1.AddonInterface

	addons, err := temp.Addons("../addons/")
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
		for _, name := range environmentConciousFilteredAddons {
			if addon.GetName() == name {
				found = true
			}
		}
		for _, name := range disabled {
			if addon.GetName() == name {
				found = true
			}
		}
		if !found {
			unhandled = append(unhandled, addon)
		}
	}

	return unhandled, nil
}
