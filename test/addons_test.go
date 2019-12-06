package test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/blang/semver"
	"gopkg.in/yaml.v2"

	"github.com/mesosphere/kubeaddons/hack/temp"
	"github.com/mesosphere/kubeaddons/pkg/api/v1beta1"
	"github.com/mesosphere/kubeaddons/pkg/test"
	"github.com/mesosphere/kubeaddons/pkg/test/cluster/kind"
)

const defaultKubernetesVersion = "1.15.6"

var addonTestingGroups = make(map[string][]string)

// -----------------------------------------------------------------------------
// Test Groups
// -----------------------------------------------------------------------------

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

func TestStorageGroup(t *testing.T) {
	if err := testgroup(t, "storage"); err != nil {
		t.Fatal(err)
	}
}

// -----------------------------------------------------------------------------
// Test Validations
// -----------------------------------------------------------------------------

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
		if !found {
			unhandled = append(unhandled, addon)
		}
	}

	return unhandled, nil
}

func init() {
	b, err := ioutil.ReadFile("groups.yaml")
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, addonTestingGroups); err != nil {
		panic(err)
	}
}
