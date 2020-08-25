SHELL := /bin/bash -euo pipefail

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
	rm -rf kubeaddons-tests
