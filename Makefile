SHELL := /bin/bash -euo pipefail
YAMLLINT := $(shell command -v yamllint)

export ADDON_TESTS_PER_ADDON_WAIT_DURATION := 10m
export ADDON_TESTS_SETUP_WAIT_DURATION := 30m
export GIT_TERMINAL_PROMPT := 1
export GO111MODULE := on
export GOPRIVATE := github.com/mesosphere/kubeaddons,github.com/mesosphere/ksphere-testing-framework
export KUBECONFIG = $(KBA_KUBECONFIG)

.DEFAULT_GOAL := test

.PHONY: elasticsearch-group-test
elasticsearch-group-test: kubeconfig set-git-ssh
	export KBA_KUBECONFIG=$(shell pwd)/kubeconfig; cd test; go test -tags experimental -timeout 60m -race -v -run TestElasticsearchGroup


kubeconfig:
ifeq (, $(shell which kind))
 $(error "No kind in $(PATH), consider installing kind")
endif
	touch kubeconfig
	kind create cluster --kubeconfig kubeconfig

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
