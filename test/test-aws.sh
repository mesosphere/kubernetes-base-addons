#!/bin/bash
set -euo pipefail

branch=${1:-master}

./scripts/setup-konvoy.sh

echo "INFO: the following test groups will be run:"
tests=$(go run -tags experimental scripts/test-wrapper.go origin ${branch} | (egrep '^Test' || true) | grep 'TestAwsGroup\|TestElasticsearchGroup')
echo ${tests}

for g in ${tests}
do
    go test -tags experimental -timeout 60m -race -v -run $g -kba-branch ${branch}
done
