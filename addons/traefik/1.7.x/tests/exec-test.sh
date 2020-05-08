#!/bin/bash
SERVER=$(kubectl -n kubeaddons get svc traefik-kubeaddons \
    -o jsonpath="{.status.loadBalancer.ingress[0].ip}{.status.loadBalancer.ingress[0].hostname}")
SECRET=$(kubectl -n default get serviceaccount kuttlaccount -o jsonpath="{ .secrets[0].name }")
TOKEN=$(kubectl -n default get secret "${SECRET}" -o go-template="{{.data.token | base64decode }}")
kubectl config set-credentials kuttl --token="${TOKEN}"
kubectl config set-context kuttl --cluster=cluster --user=kuttl
kubectl config use-context kuttl
kubectl -v 9 -s https://"${SERVER}"/testpath exec -ti testpod -- echo SUCCESS 2>.kube/kubectl.log