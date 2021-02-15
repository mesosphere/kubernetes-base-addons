MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
SHELL := /bin/bash -euo pipefail
comma := ,
empty :=
space := $(empty) $(empty)
commaspace := $(comma)$(empty)


# ------------------------------------------------------------------------------
# Configuration - Versions
# ------------------------------------------------------------------------------
GITHUB_CLI_VERSION := 1.5.0

# ------------------------------------------------------------------------------
# Configuration - Golang
# ------------------------------------------------------------------------------

GOARCH ?= $(shell go env GOARCH)
GOOS ?= $(shell go env GOOS)
GOPATH ?= $(shell go env GOPATH)
GOPRIVATE ?= "github.com/mesosphere"

ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

export GO111MODULE := on

# ------------------------------------------------------------------------------
# Configuration - Binaries
# ------------------------------------------------------------------------------
GITHUB_CLI_BIN := $(MKFILE_DIR)/bin/linux/$(GOARCH)/gh-$(GITHUB_CLI_VERSION)
RELEASE_NOTES_TOOL_BIN := $(MKFILE_DIR)/bin/$(GOOS)/$(GOARCH)/release-notes

# ------------------------------------------------------------------------------
# Configuration - Other
# ------------------------------------------------------------------------------

YAMLLINT := $(shell command -v yamllint)

export ADDON_TESTS_PER_ADDON_WAIT_DURATION := 10m
export ADDON_TESTS_SETUP_WAIT_DURATION := 30m
export GIT_TERMINAL_PROMPT := 1
export KBA_KUBECONFIG ?= $(shell mktemp --tmpdir kba-kubeconfig-XXXXXXXX)
export KUBECONFIG = $(KBA_KUBECONFIG)

.DEFAULT_GOAL := test
ADDON_SOURCES := $(wildcard addons/*/*.yaml)

ifneq (,$(filter tar (GNU tar)%, $(shell tar --version)))
WILDCARDS := --wildcards
endif

# ------------------------------------------------------------------------------
# Main
# ------------------------------------------------------------------------------

.PHONY: set-git-ssh
set-git-ssh:
ifdef DISPATCH_CI
	./scripts/ci/setup_ssh.sh
endif

# Target to run tests on Dispatch CI with KUBECONFIG from Cluster Claim Controller.
# The KUBECONFIG is set to config file in the git-clone repo of Dispatch.
.PHONY: dispatch-test
dispatch-test: set-git-ssh
	KBA_KUBECONFIG=/workspace/kba-git-src/kubeconfig ./test/dispatch-ci.sh

.PHONY: lint
lint:
	yamllint --config-file test/yamllint.yaml .

.PHONY: test
test:
	./test/run-tests.sh

kubeaddons-tests:
	git clone --depth 1 https://github.com/mesosphere/kubeaddons-tests.git --branch master --single-branch

.PHONY: test-nightly
test-nightly:
	cd test; go test -tags experimental,nightly -timeout 60m -race -v -run TestUnmarshallPrometheusMetricNames
	cd test; go test -tags experimental,nightly -timeout 60m -race -v -run TestNightlyGroup

.PHONY: ci.test-nightly
ci.test-nightly:
	# go tests
	git config --global url."https://$$GITHUB_TOKEN:@github.com/".insteadOf "https://github.com/"
	git fetch

	# docker login to get around rate limit issues, e.g. 'failed to pull image "kindest/node [...] failed with error: exit status 1'
	docker login -u $$DOCKERHUB_ROBOT_USERNAME -p $$DOCKERHUB_ROBOT_TOKEN

	cd test; ./scripts/setup-konvoy.sh

	make test-nightly

.PHONY: kind-test
kind-test: kubeaddons-tests
	make -f kubeaddons-tests/Makefile kind-test KUBEADDONS_REPO=kubernetes-base-addons

.PHONY: clean
clean:
ifneq (,$(wildcard kubeaddons-tests/Makefile))
	make -f kubeaddons-tests/Makefile clean KUBEADDONS_REPO=kubernetes-base-addons
endif
	-rm -rf kubeaddons-tests
	-rm kba-kubeconfig-*

.PHONY: dispatch-test-install-upgrade
dispatch-test-install-upgrade:
	@{ \
	echo "INFO: the following test groups will be run:" ;\
	KBA_KUBECONFIG=/workspace/kba-git-src/kubeconfig ./test/dispatch-ci.sh ;\
	cd ./test && go run -tags experimental ./scripts/test-wrapper.go ;\
	for g in $(cd ./test && go run -tags experimental ./scripts/test-wrapper.go) ; do \
	    shell cd ./test && go test -tags experimental -timeout 60m -race -v -run $g ; \
	done ;\
	}

# ------------------------------------------------------------------------------
# Release
# ------------------------------------------------------------------------------
RELEASE_LIST := $(sort $(subst $(comma), ,$(KBA_TAGS)))
RELEASE_VER := $(filter v%,$(RELEASE_LIST))

.PHONY: release.pr
release.pr: $(RELEASE_NOTES_TOOL_BIN) ADDONS.md
ifndef KBA_MILESTONE
	echo "Please set KBA_MILESTONE"
else
ifndef KBA_TAGS
	echo "Please set KBA_TAGS"
else
	echo -e "# Release Notes\n" > NEW_RELEASE_NOTES.md
	echo -e "## $(subst $(space),$(comma) ,$(RELEASE_LIST))\n" >> NEW_RELEASE_NOTES.md
	$(RELEASE_NOTES_TOOL_BIN) >> NEW_RELEASE_NOTES.md
	tail -n +3 RELEASE_NOTES.md >> NEW_RELEASE_NOTES.md
	mv NEW_RELEASE_NOTES.md RELEASE_NOTES.md
endif
endif

.PHONY: release
release: $(RELEASE_NOTES_TOOL_BIN) 
ifndef KBA_MILESTONE
	echo "Please set KBA_MILESTONE"
else
ifndef KBA_TAGS
	echo "Please set KBA_TAGS"
else
	git checkout $(KBA_MILESTONE)
	git pull
	echo $(RELEASE_LIST) | xargs -n1 echo git tag && echo git push --tags
	$(RELEASE_NOTES_TOOL_BIN) > DELETE_ME.md
	$(GITHUB_CLI_BIN) release create $(RELEASE_VER) -t $(RELEASE_VER) --target $(RELEASE_VER) --notes-file DELETE_ME.md
	rm DELETE_ME.md
endif
endif

.PHONY: make.addons.table
make.addons.table: ADDONS.md

ADDONS.md: $(ADDON_SOURCES)
	scripts/make_addon_table.sh > ADDONS.md

# ------------------------------------------------------------------------------
# Tools
# ------------------------------------------------------------------------------

.PHONY: tools
tools: $(RELEASE_NOTES_TOOL_BIN) $(GITHUB_CLI_BIN)

.PHONY: tool.release-notes
tool.release-notes: $(RELEASE_NOTES_TOOL_BIN)

.PHONY: tool.github_cli
tool.github_cli: $(GITHUB_CLI_BIN)

$(RELEASE_NOTES_TOOL_BIN): tools/cmd/release-notes/release-notes.go
	mkdir -p $(dir $@)
	cd tools/cmd/release-notes && go build -o $@ .

$(GITHUB_CLI_BIN):
	mkdir -p $(dir $@) _build
	curl -Ls https://github.com/cli/cli/releases/download/v$(GITHUB_CLI_VERSION)/gh_$(GITHUB_CLI_VERSION)_linux_$(GOARCH).tar.gz | tar xz -C _build $(WILDCARDS) --strip=2 '*/*/gh'
	mv _build/gh $@
