#!/bin/bash
# NOTE: used by teamcity and dispatch

UNAME=$(uname | tr '[:upper:]' '[:lower:]')
KONVOY_VERSION="${KONVOY_VERSION:-v1.6.1}"
if ! [ -z $1 ]
then
    KONVOY_VERSION=$1
fi

set -euo pipefail

if [[ ! -f konvoy ]] || [[ "$(./konvoy --version | awk '/Version/ {gsub("\"","",$2); gsub(",","",$2); print $2}')" == "${KONVOY_VERSION}" ]]; then
  curl --silent https://downloads.mesosphere.io/konvoy/konvoy_${KONVOY_VERSION}_${UNAME}.tar.bz2 | tar xjv --strip=1 konvoy_${KONVOY_VERSION}/konvoy
fi
