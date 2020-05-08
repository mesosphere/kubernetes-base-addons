#!/bin/bash


function deleteCluster() {
    kind delete cluster
    rm ${KUBECONFIG}
}

trap deleteCluster EXIT

kubectl kuttl test --artifacts-dir=$(ARTIFACTS) 2>&1 |tee /dev/fd/2 | go-junit-report -set-exit-code > dist/addons_test_report.xml