#!/bin/bash

PKG=bytezero_server
DATE_FOMAT=`date '+%Y%m%d-%H%M'`

TARGET_PKG=${PKG}
TAR_PKG=${PKG}_${DATE_FOMAT}.tar.gz

make

echo "tar ${TAR_PKG} begin.."
mkdir -p ${TARGET_PKG}
cp -rf bin/bytezero ${TARGET_PKG}
cp -rf public/ ${TARGET_PKG}/public
cp -rf scripts/ ${TARGET_PKG}/scripts

tar zcvf ${TAR_PKG} ${TARGET_PKG}
rm -rf ${TARGET_PKG}
echo "tar ${TAR_PKG} over..."
