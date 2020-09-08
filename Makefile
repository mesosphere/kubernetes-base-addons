SHELL := /bin/bash -euo pipefail
YAMLLINT := $(shell command -v yamllint)

export ADDON_TESTS_PER_ADDON_WAIT_DURATION := 10m
export ADDON_TESTS_SETUP_WAIT_DURATION := 30m
export GIT_TERMINAL_PROMPT := 1
export GO111MODULE := on
export GOPRIVATE := github.com/mesosphere/kubeaddons,github.com/mesosphere/ksphere-testing-framework
export KBA_KUBECONFIG ?= $(shell mktemp --tmpdir kba-kubeconfig-XXXXXXXX)
export KUBECONFIG = $(KBA_KUBECONFIG)

.DEFAULT_GOAL := test

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

.PHONY: kind-test
kind-test:
	# TODO: this is being skipped in backport branches for now, as testing requires newer versions of addons
	#       which we didn't necessarily want to backport to older Konvoy versions at the time of writing.

.PHONY: clean
clean:
ifneq (,$(wildcard kubeaddons-tests/Makefile))
	make -f kubeaddons-tests/Makefile clean KUBEADDONS_REPO=kubernetes-base-addons
endif
	-rm -rf kubeaddons-tests
	-rm kba-kubeconfig-*
