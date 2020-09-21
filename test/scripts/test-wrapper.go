package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/blang/semver"
	testgroups "github.com/mesosphere/ksphere-testing-framework/pkg/groups"
	"github.com/mesosphere/kubeaddons/pkg/catalog"
	"github.com/mesosphere/kubeaddons/pkg/constants"
	"github.com/mesosphere/kubeaddons/pkg/repositories"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
	testutils "github.com/mesosphere/kubeaddons/test/utils"
)

const (
	upstreamRemote = "origin"
	upstreamBranch = "master"
)

type addonName string

var re = regexp.MustCompile(`^addons/([a-zA-Z-]+)/?`)

func main() {
	modifiedAddons, err := getModifiedAddons()
	if err != nil {
		panic(err)
	}

	r, err := local.NewRepository("local", "../addons/")
	if err != nil {
		panic(err)
	}

	if err := ensureModifiedAddonsHaveUpdatedRevisions(modifiedAddons, r); err != nil {
		panic(err)
	}

	c, err := catalog.NewCatalog(r)
	if err != nil {
		panic(err)
	}

	groups, err := testgroups.AddonsForGroupsFile("groups.yaml", c)
	if err != nil {
		panic(err)
	}

	atLeastOneGroupNeedsTesting := false
	for group, addons := range groups {
		included := false
		for _, addon := range addons {
			for _, addonName := range modifiedAddons {
				if addon.GetName() == string(addonName) {
					included = true
					atLeastOneGroupNeedsTesting = true
				}
			}
		}

		if included {
			fmt.Printf("Test%sGroup\n", strings.Title(string(group)))
		}
	}

	if !atLeastOneGroupNeedsTesting {
		for group := range groups {
			fmt.Printf("Test%sGroup\n", strings.Title(string(group)))
		}
	}
}

func getModifiedAddons() ([]addonName, error) {
	addonsModifiedMap := make(map[addonName]struct{})
	stdout := new(bytes.Buffer)
	cmd := exec.Command("git", "diff", "origin/master", "--name-only")
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	for _, line := range strings.Split(stdout.String(), "\n") {
		submatches := re.FindStringSubmatch(line)
		if submatches != nil {
			addonsModifiedMap[addonName(submatches[1])] = struct{}{}
		}
	}

	addonsModified := make([]addonName, 0, len(addonsModifiedMap))
	for name := range addonsModifiedMap {
		addonsModified = append(addonsModified, name)
	}

	return addonsModified, nil
}

func ensureModifiedAddonsHaveUpdatedRevisions(namesOfModifiedAddons []addonName, repo repositories.Repository) error {
	for _, addonName := range namesOfModifiedAddons {
		fmt.Printf("INFO: ensuring revision was updated for modified addon %s\n", addonName)

		modifiedAddonRevisions, err := repo.GetAddon(string(addonName))
		if err != nil {
			return err
		}

		modifiedAddon, err := modifiedAddonRevisions.Latest()
		if err != nil {
			return err
		}

		upstreamAddon, err := testutils.GetLatestAddonRevisionFromLocalRepoBranch("../", upstreamRemote, upstreamBranch, string(addonName))
		if err != nil {
			if strings.Contains(err.Error(), "directory not found") {
				fmt.Printf("%s is a new addon, revision check skipped", addonName)
				continue
			}
			return err
		}

		modifiedVersion := semver.MustParse(strings.TrimPrefix(modifiedAddon.GetAnnotations()[constants.AddonRevisionAnnotation], "v"))
		upstreamVersion := semver.MustParse(strings.TrimPrefix(upstreamAddon.GetAnnotations()[constants.AddonRevisionAnnotation], "v"))

		if modifiedVersion.LE(upstreamVersion) {
			return fmt.Errorf("the revision for addons %s was not properly updated (current: %s, previous from branch %s: %s). Please update the revision for any addons which you modify (see CONTRIBUTING.md)", addonName, modifiedVersion, upstreamBranch, upstreamVersion)
		}

		fmt.Printf("INFO: addon %s has an updated revision %s (upstream branch %s: %s)\n", addonName, modifiedVersion, upstreamBranch, upstreamVersion)
	}

	return nil
}
