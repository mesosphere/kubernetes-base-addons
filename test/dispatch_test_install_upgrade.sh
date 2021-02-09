#!/bin/bash
set -euo pipefail

branch=${1:-master}

./scripts/setup-konvoy.sh

# in dispatch we should be able to create these resources, for manual testing this will not work
if [[ ! -z "${KBA_KUBECONFIG}" ]]; then
    $(kubectl --kubeconfig ${KBA_KUBECONFIG} create namespace cert-manager) || true
    $(kubectl --kubeconfig ${KBA_KUBECONFIG} create secret tls kubernetes-root-ca --namespace=cert-manager --cert=/etc/kubernetes/pki/ca.crt --key=/etc/kubernetes/pki/ca.key --dry-run -o yaml  | kubectl --kubeconfig ${KBA_KUBECONFIG} apply -f - ) || true
fi


echo "INFO: the following test groups will be run:"
tests=$(go run -tags experimental scripts/test-wrapper.go origin ${branch} | (egrep '^Test' || true))
echo ${tests}

for g in ${tests}
do
    go test -tags experimental -timeout 60m -race -v -run $g -kba-branch ${branch}
done
