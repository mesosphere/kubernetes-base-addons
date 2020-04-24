#!/bin/bash
# NOTE: used by teamcity and dispatch

KONVOY_VERSION="v1.4.1"
if ! [ -z $1 ]
then
    KONVOY_VERSION=$1
fi

set -euo pipefail

curl --silent https://downloads.mesosphere.io/konvoy/konvoy_${KONVOY_VERSION}_linux.tar.bz2 -o - | tar --extract -O -jf - konvoy_${KONVOY_VERSION}/konvoy > konvoy

chmod +x konvoy
