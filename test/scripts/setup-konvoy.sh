#!/bin/bash
# NOTE: used by teamcity and dispatch

KONVOY_VERSION="v1.5.0-beta.4"
if ! [ -z $1 ]
then
    KONVOY_VERSION=$1
fi

set -euo pipefail

curl --silent https://downloads.mesosphere.io/konvoy/konvoy_${KONVOY_VERSION}_linux.tar.bz2 | tar xjv --strip=1 konvoy_${KONVOY_VERSION}/konvoy
