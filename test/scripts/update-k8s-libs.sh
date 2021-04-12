#!/bin/env bash

if [[ -z "$1" ]]; then
    echo "Usage: $0 <version>"
    echo
    echo "Example: $0 v0.19.4"
    exit 1
fi

KUBE_LIB_VERSION="$1"

go get -t \
    k8s.io/apiextensions-apiserver@${KUBE_LIB_VERSION} \
    k8s.io/apimachinery@${KUBE_LIB_VERSION} \
    k8s.io/client-go@${KUBE_LIB_VERSION} \
    k8s.io/api/core/v1@${KUBE_LIB_VERSION} \
    k8s.io/apiextensions-apiserver@${KUBE_LIB_VERSION} \
    k8s.io/apimachinery@${KUBE_LIB_VERSION} \
    k8s.io/apiserver@${KUBE_LIB_VERSION} \
    k8s.io/cli-runtime@${KUBE_LIB_VERSION} \
    k8s.io/client-go@${KUBE_LIB_VERSION} \
    k8s.io/cloud-provider@${KUBE_LIB_VERSION} \
    k8s.io/cluster-bootstrap@${KUBE_LIB_VERSION} \
    k8s.io/code-generator/cmd/client-gen@${KUBE_LIB_VERSION} \
    k8s.io/component-base@${KUBE_LIB_VERSION} \
    k8s.io/cri-api@${KUBE_LIB_VERSION} \
    k8s.io/csi-translation-lib@${KUBE_LIB_VERSION} \
    k8s.io/kube-aggregator@${KUBE_LIB_VERSION} \
    k8s.io/kube-controller-manager@${KUBE_LIB_VERSION} \
    k8s.io/kube-proxy@${KUBE_LIB_VERSION} \
    k8s.io/kube-scheduler@${KUBE_LIB_VERSION} \
    k8s.io/kubectl@${KUBE_LIB_VERSION} \
    k8s.io/kubelet@${KUBE_LIB_VERSION} \
    k8s.io/legacy-cloud-providers@${KUBE_LIB_VERSION} \
    k8s.io/metrics@${KUBE_LIB_VERSION} \
    k8s.io/sample-apiserver@${KUBE_LIB_VERSION}

go mod tidy
