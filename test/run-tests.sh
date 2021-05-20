#!/bin/bash

# This script is invoked by make test, and exercises integration tests.

set -euxo pipefail

# Tests are exercised from `test` directory. The current directory of the script.
cd "$(dirname "$0")"

echo "Setup Konvoy"
source ./scripts/setup-konvoy.sh v1.7.2

echo "git fetch branches."
git fetch

echo "INFO: the following test groups will be run:"
TESTS=($(go run -tags experimental scripts/test-wrapper.go))
echo ${TESTS[*]}

for g in ${TESTS[*]}
do
    go test -tags experimental -timeout 60m -race -v -run "$g"
done
