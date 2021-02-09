SHELL := /bin/bash -euo pipefail
YAMLLINT := $(shell command -v yamllint)

export ADDON_TESTS_PER_ADDON_WAIT_DURATION := 10m
export ADDON_TESTS_SETUP_WAIT_DURATION := 30m
export GIT_TERMINAL_PROMPT := 1
export GO111MODULE := on
export GOPRIVATE := github.com/mesosphere/kubeaddons,github.com/mesosphere/ksphere-testing-framework
export KBA_KUBECONFIG ?= /workspace/kba-git-src/kubeconfig
#export KUBECONFIG = $(KBA_KUBECONFIG)
export KBA_BRANCH ?= $(shell git branch | grep -v detached | awk '$$1=="*"{print $$2}')

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
	./test/dispatch-ci.sh

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

.PHONY: release
release: make.addons.table

.PHONY: make.addons.table
make.addons.table:
<<<<<<< HEAD
	scripts/make_addon_table.sh > ADDONS.md
=======
	scripts/make_addon_table.sh > ADDONS.md

.PHONY: dispatch-test-install-upgrade
<<<<<<< HEAD
dispatch-test-install-upgrade: set-git-ssh
	./test/scripts/setup-konvoy.sh

	KBA_KUBECONFIG=/workspace/kba-git-src/kubeconfig ./test/dispatch-ci.sh
	echo "INFO: the following test groups will be run:"
	cd ./test && go run -tags experimental ./scripts/test-wrapper.go

	for g in $(shell cd ./test && go run -tags experimental ./scripts/test-wrapper.go) ; do \
	    shell cd ./test && go test -tags experimental -timeout 60m -race -v -run $g ; \
	done
>>>>>>> e3a2c47... chore: split and parallize install and upgrade jobs.
=======
dispatch-test-install-upgrade:
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
	@{ \
	echo "INFO: the following test groups will be run:" ;\
	KBA_KUBECONFIG=/workspace/kba-git-src/kubeconfig ./test/dispatch-ci.sh ;\
	cd ./test && go run -tags experimental ./scripts/test-wrapper.go ;\
	for g in $(cd ./test && go run -tags experimental ./scripts/test-wrapper.go) ; do \
		cd ./test && go test -tags experimental -timeout 60m -race -v -run $g ; \
	done ;\
	}
>>>>>>> 56c5dac... Fixes makefile target for Dispatch
=======
	unset KBA_KUBECONFIG && cd ./test/ && ./test_install_upgrade.sh $(KBA_BRANCH)
>>>>>>> de4b996... move test to local bash script
=======
	cd ./test/ && KBA_KUBECONFIG=/workspace/kba-git-src/kubeconfig ./test_install_upgrade.sh $(KBA_BRANCH)
>>>>>>> c9e0a5f... set kubeconfig to avoid an OOM
=======
=======
	ls /workspace/kba-git-src/
	cat /workspace/kba-git-src/kubeconfig
>>>>>>> 0e2d737... Debug and ignore dispatchfile
=======
>>>>>>> 379ed74... create resources on dispath kind cluster
	cd test; ./dispatch_test_install_upgrade.sh $(KBA_BRANCH)
>>>>>>> 3ca04a7... multiple fixes to the makefile and dispatch
