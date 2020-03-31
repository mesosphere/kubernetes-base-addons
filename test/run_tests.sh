#!/bin/bash 
set -euxo pipefail

git fetch

echo "INFO: the following test groups will be run:"
go run -tags experimental scripts/test-wrapper.go

for g in $(go run scripts/test-wrapper.go); do
	go test -timeout 30m -race -v -run $g;
done
