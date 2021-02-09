#!/bin/bash
set -euo pipefail

branch=${1:-master}

./scripts/setup-konvoy.sh

# in dispatch we should be able to create these resources, for manual testing this will not work
if [[ ! -z "${KBA_KUBECONFIG}" ]]; then
    $(kubectl --kubeconfig ${KBA_KUBECONFIG} create namespace cert-manager || true)
    KIND_POD_NAME=$(kubectl get clusterclaim $CLAIM_NAME | awk '{print $2}')
    kubectl exec  -ti -n dispatch $KIND_POD_NAME -- bash -c "docker ps -aq | xargs -I{} docker exec {} cat /etc/kubernetes/pki/ca.crt"  >> ca.crt
    kubectl exec  -ti -n dispatch $KIND_POD_NAME -- bash -c "docker ps -aq | xargs -I{} docker exec {} cat /etc/kubernetes/pki/ca.key"  >> ca.key
    $(kubectl --kubeconfig ${KBA_KUBECONFIG} create secret tls kubernetes-root-ca --namespace=cert-manager --cert=ca.crt --key=ca.key --dry-run -o yaml  | kubectl --kubeconfig ${KBA_KUBECONFIG} apply -f - ) || true
fi


echo "INFO: the following test groups will be run:"
tests=$(go run -tags experimental scripts/test-wrapper.go origin ${branch} | (egrep '^Test' || true))
echo ${tests}

for g in ${tests}
do
    go test -tags experimental -timeout 60m -race -v -run $g -kba-branch ${branch}
done
