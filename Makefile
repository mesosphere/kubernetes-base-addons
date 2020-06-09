SHELL := /bin/bash -euo pipefail

export GO111MODULE := on
export ADDON_TESTS_PER_ADDON_WAIT_DURATION := 10m
export GIT_TERMINAL_PROMPT := 1
export ADDON_TESTS_SETUP_WAIT_DURATION := 30m
export GOPRIVATE := github.com/mesosphere/kubeaddons

.DEFAULT_GOAL := test

.PHONY: set-git-ssh
set-git-ssh:
ifdef DISPATCH_CI
	./scripts/ci/setup_ssh.sh
endif

# Target to run restricted set of tests on Dispatch CI.
.PHONY: dispatch-test
dispatch-test: set-git-ssh
	export KUBECONFIG=`pwd`/kubeconfig
	./test/dispatch-ci.sh

.PHONY: test
test:
	./test/run-tests.sh
