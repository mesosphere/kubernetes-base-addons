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

.PHONY: test
test: set-git-ssh
	cd test && git fetch; \
	for g in $(shell cd test && go run scripts/test-wrapper.go); do \
		go test -timeout 30m -race -v -run $$g; \
	done
