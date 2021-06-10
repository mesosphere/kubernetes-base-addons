#!/bin/bash
set -xeuo pipefail

if [[ ${DISPATCH_BUILD_NAME#pr-} =~ ^[0-9]+$ ]]; then
  branch=$(${GITHUB_CLI_BIN} api --method GET repos/:owner/:repo/pulls/${DISPATCH_BUILD_NAME#pr-} -q .base.ref)
else
  branch=master
fi

# if we're in dispatch - create the secrets on the container running in the pod with exec
if [[ -n "${CLAIM_NAME}" ]]; then
    kubectl --kubeconfig ${KBA_KUBECONFIG} create namespace cert-manager || true
    KIND_POD_NAME=$(kubectl get clusterclaim $CLAIM_NAME | tail +2 | awk '{print $2}')
    kubectl exec -n dispatch $KIND_POD_NAME -- bash -c "docker ps -aq | xargs -I{} docker exec {} cat /etc/kubernetes/pki/ca.crt"  >> ca.crt
    kubectl exec -n dispatch $KIND_POD_NAME -- bash -c "docker ps -aq | xargs -I{} docker exec {} cat /etc/kubernetes/pki/ca.key"  >> ca.key
    kubectl --kubeconfig ${KBA_KUBECONFIG} create secret tls kubernetes-root-ca --namespace=cert-manager --cert=ca.crt --key=ca.key --dry-run -o yaml  | kubectl --kubeconfig ${KBA_KUBECONFIG} apply -f -  || true
fi


echo "INFO: the following test groups will be run:"
tests=$(go run -tags experimental scripts/test-wrapper.go origin ${branch} | (grep -E '^Test' || true) | grep -Ev 'TestAwsGroup\|TestElasticsearchGroup\|TestIstioGroup')
echo ${tests}

for g in ${tests}
do
    go test -tags experimental -timeout 60m -race -v -run $g -kba-branch ${branch}
done
