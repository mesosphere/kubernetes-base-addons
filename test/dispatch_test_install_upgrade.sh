#!/bin/bash
set -euo pipefail

branch=${1:-release/3}
git fetch origin ${branch} 

# if we're in dispatch - create the secrets on the container running in the pod with exec
if [[ ! -z "${CLAIM_NAME}" ]]; then
    kubectl --kubeconfig ${KBA_KUBECONFIG} create namespace cert-manager || true
    KIND_POD_NAME=$(kubectl get clusterclaim $CLAIM_NAME | tail +2 | awk '{print $2}')
    kubectl exec -n dispatch $KIND_POD_NAME -- bash -c "docker ps -aq | xargs -I{} docker exec {} cat /etc/kubernetes/pki/ca.crt"  >> ca.crt
    kubectl exec -n dispatch $KIND_POD_NAME -- bash -c "docker ps -aq | xargs -I{} docker exec {} cat /etc/kubernetes/pki/ca.key"  >> ca.key
    kubectl --kubeconfig ${KBA_KUBECONFIG} create secret tls kubernetes-root-ca --namespace=cert-manager --cert=ca.crt --key=ca.key --dry-run -o yaml  | kubectl --kubeconfig ${KBA_KUBECONFIG} apply -f -  || true
fi

echo "INFO: the following test groups will be run against branch ${branch}:"
tests=$(go run -tags experimental scripts/test-wrapper.go origin ${branch} | (egrep '^Test' || true) | grep -v 'TestAwsGroup\|TestElasticsearchGroup\|TestIstioGroup')
echo ${tests}

for g in ${tests}
do
    go test -tags experimental -timeout 60m -race -v -run $g -kba-branch ${branch}
done
