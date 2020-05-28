SHELL := /bin/bash -euo pipefail
RELEASE_NEXT_VER := $(shell git describe --tags --always origin/testing | sed 's/testing-//' | awk -F. '{print $1"."$2".0"}')
SOAK_NEXT_VER := $(shell git describe --tags --always origin/testing | sed 's/testing-//' | awk -F. '{print $1"."($2+1)".0"}')

export GO111MODULE := on
export ADDON_TESTS_PER_ADDON_WAIT_DURATION := 10m
export GIT_TERMINAL_PROMPT := 1
export ADDON_TESTS_SETUP_WAIT_DURATION := 30m
export GOPRIVATE := github.com/mesosphere/kubeaddons,github.com/mesosphere/ksphere-testing-framework

.DEFAULT_GOAL := test

.PHONY: set-git-ssh
set-git-ssh:
ifdef DISPATCH_CI
	./scripts/ci/setup_ssh.sh
endif

# Target to run restricted set of tests on Dispatch CI.
.PHONY: dispatch-test
dispatch-test: set-git-ssh
	./test/dispatch-ci.sh

.PHONY: test
test:
	./test/run-tests.sh

.PHONY: release
release:
	if [ -z '${GITHUB_TOKEN}' ]; then echo 'Environment variable GITHUB_TOKEN not set' && exit 1; fi
	git checkout stable
	git fetch --all
	git reset --hard origin stable
	git checkout -b stable-$(RELEASE_NEXT_VER)
	git merge -s recursive -X theirs origin/testing
	rm /tmp/rn || true
	release-notes --start-sha $$(git rev-parse stable) --end-sha $$(git rev-parse HEAD) --github-org mesosphere --github-repo kubernetes-base-addons --required-author "" --format json --output /tmp/rn
	cat <(echo -e "## stable-1.15-$(RELEASE_NEXT_VER), stable-1.16-$(RELEASE_NEXT_VER)\n") \
		<(jq -r '"* " + .[].markdown' /tmp/rn) \
		<(echo) \
		<(cat RELEASE_NOTES.md) > RELEASE_NOTES.tmp
	mv RELEASE_NOTES.tmp RELEASE_NOTES.md
	git add RELEASE_NOTES.md
	git commit -m "docs: add release notes for stable-$(RELEASE_NEXT_VER)"
	git push -u origin stable-$(RELEASE_NEXT_VER)
	curl -u x:${GITHUB_TOKEN} -X POST \
		--data '{"title": "release: stable-$(RELEASE_NEXT_VER)", "head": "stable-$(RELEASE_NEXT_VER)", "base": "stable", "body": "Release of stable-$(RELEASE_NEXT_VER)"}' \
		"https://api.github.com/repos/mesosphere/kubernetes-base-addons/pulls"

.PHONY: testing-branch
testing-branch:
	if [ -z '${GITHUB_TOKEN}' ]; then echo 'Environment variable GITHUB_TOKEN not set' && exit 1; fi
	if ! git merge-base --is-ancestor $$(git rev-parse origin/stable) origin/master; then echo 'stable must be merged into master before creating a testing branch' && exit 1; fi 
	git checkout testing
	git fetch --all
	git reset --hard origin testing
	git merge -s recursive -X theirs origin/master
	git tag testing-$(SOAK_NEXT_VER)
	git push -u origin testing
