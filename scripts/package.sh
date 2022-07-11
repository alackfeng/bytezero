#!/bin/bash

PKG=bytezero_server
DATE_FOMAT=`date '+%Y%m%d-%H%M'`

TARGET_PKG=${PKG}
TAR_PKG=${PKG}_${DATE_FOMAT}.tar.gz

make

echo "tar ${TAR_PKG} begin.."
mkdir -p ./install/${TARGET_PKG}/bin
cp -rf bin/bytezero ./install/${TARGET_PKG}/bin/
cp -rf public/ ./install/${TARGET_PKG}/public
cp -rf scripts/ ./install/${TARGET_PKG}/scripts

cd ./install
tar zcvf ${TAR_PKG} ${TARGET_PKG}
cd ..
echo "tar ./install/${TAR_PKG} over..."
rm -rf ./install/${TARGET_PKG}
echo "sz ./install/${TAR_PKG}"
