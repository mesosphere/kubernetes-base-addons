SHELL := /bin/bash -euo pipefail
YAMLLINT := $(shell command -v yamllint)

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

# Target to run tests on Dispatch CI with KUBECONFIG from Cluster Claim Controller.
# The KUBECONFIG is set to config file in the git-clone repo of Dispatch.
.PHONY: dispatch-test
dispatch-test: set-git-ssh
	KUBECONFIG=/workspace/kba-git-src/kubeconfig ./test/dispatch-ci.sh

.PHONY: lint
lint:
	yamllint --config-file test/yamllint.yaml .

.PHONY: test
test:
	./test/run-tests.sh
