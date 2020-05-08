SHELL := /bin/bash -euo pipefail

export GO111MODULE := on
export ADDON_TESTS_PER_ADDON_WAIT_DURATION := 10m
export GIT_TERMINAL_PROMPT := 1
export ADDON_TESTS_SETUP_WAIT_DURATION := 30m
export GOPRIVATE := github.com/mesosphere/kubeaddons
export KUBECONFIG := kubeconfig

KUTTL_VERSION=0.1.0
KIND_VERSION=0.8.1
KUBERNETES_VERSION ?= 1.17.5


OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
MACHINE=$(shell uname -m)
KIND_MACHINE=$(shell uname -m)
ifeq "$(KIND_MACHINE)" "x86_64"
  KIND_MACHINE=amd64
endif

export PATH := $(shell pwd)/bin/:$(PATH)

ARTIFACTS=dist

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

bin/.placeholder:
	mkdir -p bin/
	touch bin/.placeholder

bin/kind_$(KIND_VERSION): bin/.placeholder
	curl -Lso bin/kind_$(KIND_VERSION) https://github.com/kubernetes-sigs/kind/releases/download/v$(KIND_VERSION)/kind-$(OS)-$(KIND_MACHINE)
	chmod +x bin/kind_$(KIND_VERSION)

.PHONY: bin/kind
bin/kind: bin/kind_$(KIND_VERSION)
	@ln -sf ./kind_$(KIND_VERSION) bin/kind

bin/kubectl-kuttl_$(KUTTL_VERSION): bin/.placeholder
	curl -Lso bin/kubectl-kuttl_$(KUTTL_VERSION) https://github.com/kudobuilder/kuttl/releases/download/v$(KUTTL_VERSION)/kubectl-kuttl_$(KUTTL_VERSION)_$(OS)_$(MACHINE)
	chmod +x bin/kubectl-kuttl_$(KUTTL_VERSION)

.PHONY: bin/kubectl-kuttl
bin/kubectl-kuttl: bin/kubectl-kuttl_$(KUTTL_VERSION)
	@ln -sf ./kubectl-kuttl_$(KUTTL_VERSION) bin/kubectl-kuttl

bin/go-junit-report:
	GOBIN=${PWD}/bin/ go get -u github.com/jstemmer/go-junit-report

.PHONY: install-bin
install-bin: bin/kind bin/kubectl-kuttl bin/go-junit-report

.PHONY: create-kind-cluster
create-kind-cluster: install-bin $(KUBECONFIG)

.PHONY: delete-kind-cluster
delete-kind-cluster:
	kind delete cluster
	rm $(KUBECONFIG)

$(KUBECONFIG):
	kind create cluster --wait 10s --config=test/kind/kubernetes-$(KUBERNETES_VERSION).yaml
	mkdir -p $(ARTIFACTS)
	kubectl kuttl test --artifacts-dir=$(ARTIFACTS) test/kind/init 2>&1 |tee /dev/fd/2 | go-junit-report -set-exit-code > dist/addons_setup_report.xml

.PHONY: kind-kuttl-test
kind-kuttl-test: create-kind-cluster
	mkdir -p $(ARTIFACTS)
	scripts/kind-kuttl-test.sh
