#!/bin/bash
set -eo pipefail

branch=${1:-master}

echo "INFO: the following test groups will be run:"
tests=$(go run -tags experimental scripts/test-wrapper.go origin "${branch}" | grep '^Test'  | grep 'TestAwsGroup\|TestElasticsearchGroup\|TestIstioGroup') && echo ${tests}

pids=()
for g in ${tests}
do
    go test -tags experimental -timeout 60m -race -v -run "${g}" -kba-branch "${branch}" > "${g}.log" 2>&1 &
    pids+=($!)
done

exits=()
EXIT_CODE=0
for pid in ${pids[*]}; do
  wait "${pid}" || EXIT_CODE=$?
  exits+=("${EXIT_CODE}")
  if [ "${DISPATCH_CI}" = "true" ]; then
    cat Test*
    cp Test* "${ARTIFACTS_DIR}"
  fi
done

echo "${exits[*]}"

for e in ${exits[*]}; do
  if [ "${e}" -ne 0 ]; then
    exit "${e}"
  fi
done
