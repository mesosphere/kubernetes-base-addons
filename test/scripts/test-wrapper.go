package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mesosphere/ksphere-testing-framework/pkg/experimental"
	"github.com/mesosphere/kubeaddons/pkg/catalog"
	"github.com/mesosphere/kubeaddons/pkg/repositories/local"
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

	c, err := catalog.NewCatalog(r)
	if err != nil {
		panic(err)
	}

	groups, err := experimental.AddonsForGroupsFile("groups.yaml", c)
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
