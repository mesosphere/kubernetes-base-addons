#!/bin/bash

JP='{.status.addresses[?(.type == "InternalIP")].address}'
NODEIP=$(kubectl get node kind-control-plane -o jsonpath="${JP}")
SLASH24="${NODEIP%.*}"
kubectl -n kubeaddons create configmap metallb-config --from-literal=config='{"address-pools":[{"name":"default","protocol":"layer2","addresses":["'"${SLASH24}"'.200-'"${SLASH24}"'.250"]}]}'
