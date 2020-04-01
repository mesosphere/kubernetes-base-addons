#!/bin/bash
set -euxo pipefail

git fetch

echo "INFO: the following test groups will be run:"
#go run -tags experimental scripts/test-wrapper.go

go test -tags experimental -timeout 60m -race -v -run TestGeneralGroup

#for g in $(go run -tags experimental scripts/test-wrapper.go)
#do
#    go test -tags experimental -timeout 60m -race -v -run $g
#done
