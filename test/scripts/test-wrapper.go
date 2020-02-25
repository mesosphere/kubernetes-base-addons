package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type groupName string
type addonName string

var re = regexp.MustCompile(`^addons/([a-z]+)/?`)

func main() {
	modifiedAddons, err := getModifiedAddons()
	if err != nil {
		panic(err)
	}

	testGroups, err := getGroupsToTest(modifiedAddons)
	if err != nil {
		panic(err)
	}

	for _, group := range testGroups {
		fmt.Printf("Test%sGroup\n", strings.Title(string(group)))
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

func getGroupsToTest(modifiedAddons []addonName) ([]groupName, error) {
	b, err := ioutil.ReadFile("groups.yaml")
	if err != nil {
		return nil, err
	}

	g := make(map[groupName][]addonName)
	if err := yaml.Unmarshal(b, &g); err != nil {
		return nil, err
	}

	testGroups := make([]groupName, 0)
	// if no Addon has been modified, return all existing groups
	if len(modifiedAddons) == 0 {
		for group, _ := range g {
			testGroups = append(testGroups, group)
		}
		return testGroups, nil
	}

	for _, modifiedAddonName := range modifiedAddons {
		for group, addons := range g {
			for _, name := range addons {
				if name == modifiedAddonName {
					exists := false
					for _, existingGroup := range testGroups {
						if group == existingGroup {
							exists = true
						}
					}
					if !exists {
						testGroups = append(testGroups, group)
					}
				}
			}
		}
	}

	if len(testGroups) < 1 {
		return nil, fmt.Errorf("error: there were no testGroups to test")
	}

	return testGroups, nil
}
