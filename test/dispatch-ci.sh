#!/bin/bash

# This script is exercised from Dispatch.

set -euxo pipefail

# Tests are exercised from `test` directory. The current directory of the script.
cd "$(dirname "$0")"

echo "Setup Konvoy"
source ./scripts/setup-konvoy.sh v1.5.0-beta.4

echo "git fetch branches."
git fetch

echo "INFO: Exercise TestAwsGroup "
go test -tags experimental -timeout 60m -race -v -run TestAwsGroup

# Fix kind cluster issues on Dispatch CI before enabling complete test-suite.

# Once these issues are fixed. This target and the dispatch specific test script
# can be removed and make test can be exercised directly on dispatch.
#
# TODO: (D2IQ-66356) - Fix the kind cluster issues on dispatch
#
# echo "INFO: the following test groups will be run:"
# go run -tags experimental scripts/test-wrapper.go
#
# for g in $(go run -tags experimental scripts/test-wrapper.go)
# do
#     go test -tags experimental -timeout 60m -race -v -run $g
# done
