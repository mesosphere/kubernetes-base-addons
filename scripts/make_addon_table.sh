#!/bin/bash
set -euo pipefail

gather_data(){
    for addon in addons/*/*.yaml; do
        KIND=$(yq -r '.kind' "${addon}")
        NAME=$(yq -r '.metadata.name' "${addon}")
        NAMESPACE=$(yq -r '.metadata.namespace' "${addon}" | sed 's@^null$@ @')
        VERSION=$(yq -r '.metadata.annotations["catalog.kubeaddons.mesosphere.io/addon-revision"]' "${addon}")
        HELM_CHART_VERSION=$(yq -r '.spec.chartReference.version' "${addon}")
        HELM_REPO=$(yq -r '.spec.chartReference.repo' "${addon}")
        HELM_CHART=$(yq -r '.spec.chartReference.chart' "${addon}")
        APP_VERSION=$(yq -r '.metadata.annotations["appversion.kubeaddons.mesosphere.io/'${NAME}'"]' "${addon}")

        echo "|${NAMESPACE}|${NAME}|${VERSION}|${APP_VERSION}|${KIND}|"
    done
}
echo "|Namespace|Name|Version|App Version|Kind|"
echo "|---------|----|-------|-----------|----|"

gather_data | sort -t'|' -k1,2