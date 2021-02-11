#!/bin/bash
set -euo pipefail

branch=${1:-master}

echo "INFO: the following test groups will be run:"
tests=$(go run -tags experimental scripts/test-wrapper.go origin ${branch} | grep '^Test'  | grep 'TestAwsGroup\|TestElasticsearchGroup\|TestIstioGroup')
echo ${tests}

pids=()
for g in ${tests}
do
    go test -tags experimental -timeout 60m -race -v -run $g -kba-branch ${branch} > $g 2>&1 &
    pids+=($!)
done

exits=()
for pid in ${pids[*]}; do
  wait $pid || EXIT_CODE=$?
  exits+=($EXIT_CODE)
  if [ $DISPATCH_CI = "true" ]; then
    cp Test* $ARTIFACTS_DIR
  fi
done

echo ${exits[*]}

for e in ${exits[*]}; do
  if [ ${e} -ne 0 ]; then
    exit $e
  fi
done
