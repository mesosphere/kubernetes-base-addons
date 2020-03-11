#!/bin/bash
# NOTE: used by teamcity

KONVOY_VERSION="v1.4.1"
if ! [ -z $1 ]
then
    KONVOY_VERSION=$1
fi

set -euo pipefail

wget -q https://downloads.mesosphere.io/konvoy/konvoy_${KONVOY_VERSION}_linux.tar.bz2 -O - | tar --extract -O -jf - konvoy_${KONVOY_VERSION}/konvoy > konvoy

chmod +x konvoy
