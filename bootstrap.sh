#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This code is based on code written by the Hyperledger Fabric community.
# Original code can be found here: https://github.com/hyperledger/fabric/blob/release-1.1/scripts/bootstrap.sh
#

# if version not passed in, default to latest released version
export VERSION=${1:-1.1.0}
# if ca version not passed in, default to latest released version
export CA_VERSION=${2:-$VERSION}
export ARCH=$(echo "$(uname -s|tr '[:upper:]' '[:lower:]'|sed 's/mingw64_nt.*/windows/')-$(uname -m | sed 's/x86_64/amd64/g')" | awk '{print tolower($0)}')
#Set MARCH variable i.e ppc64le,s390x,x86_64,i386
MARCH=`uname -m`

dockerFabricPull() {
  local FABRIC_TAG=$1
  for IMAGES in peer orderer tools; do
      echo "==> FABRIC IMAGE: $IMAGES"
      echo
      docker pull hyperledger/fabric-$IMAGES:$FABRIC_TAG
      docker tag hyperledger/fabric-$IMAGES:$FABRIC_TAG hyperledger/fabric-$IMAGES
  done
}

dockerCaPull() {
      local CA_TAG=$1
      echo "==> FABRIC CA IMAGE"
      echo
      docker pull hyperledger/fabric-ca:$CA_TAG
      docker tag hyperledger/fabric-ca:$CA_TAG hyperledger/fabric-ca
}

: ${CA_TAG:="$MARCH-$CA_VERSION"}
: ${FABRIC_TAG:="$MARCH-$VERSION"}

echo "===> Downloading platform specific fabric binaries"
curl https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric/hyperledger-fabric/${ARCH}-${VERSION}/hyperledger-fabric-${ARCH}-${VERSION}.tar.gz | tar xz

echo "===> Downloading platform specific fabric-ca-client binary"
curl https://nexus.hyperledger.org/content/repositories/releases/org/hyperledger/fabric-ca/hyperledger-fabric-ca/${ARCH}-${VERSION}/hyperledger-fabric-ca-${ARCH}-${VERSION}.tar.gz | tar xz
if [ $? != 0 ]; then
     echo
     echo "------> $VERSION fabric-ca-client binary is not available to download  (Avaialble from 1.1.0-rc1) <----"
     echo
fi
echo "===> Pulling fabric Images"
dockerFabricPull ${FABRIC_TAG}

echo "===> Pulling fabric ca Image"
dockerCaPull ${CA_TAG}

echo
echo "===> List out hyperledger docker images"
docker images | grep hyperledger*