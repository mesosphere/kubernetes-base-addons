#!/bin/bash

# This script is exercised from Dispatch.

set -euxo pipefail

# Tests are exercised from `test` directory. The current directory of the script.
cd "$(dirname "$0")"

echo "Setup Konvoy"
source ./scripts/setup-konvoy.sh v1.6.1

echo "git fetch branches."
git fetch

echo "Run a subset of tests on dispatch"

# Full Testsuite is Blocked on https://jira.d2iq.com/browse/D2IQ-70406
# Dispatch Staging Cluster cannot run full KBA test suite. Results in Pipeline timeout

DISPATCH_SUBSET_REGEX='TestDisabledGroup\|TestGeneralGroup\|TestAwsGroup'

go test -tags experimental -timeout 60m -race -v -run TestIstioGroup
