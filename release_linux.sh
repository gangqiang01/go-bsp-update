#!/bin/bash
ClientVersion=$(cat common/constants/default.go  | grep ClientVersion | sed 's/\"//g')
ClientVersion=${ClientVersion#*=}

V=$(echo $ClientVersion |awk '{print $1}' |sed 's/ //g' | sed 's/\t//g')
V=${V:-Beta}
Version=$(echo $ClientVersion |awk '{print $2}' |sed 's/ //g' | sed 's/\t//g')
#Version=${Version:-0.2.0}

TARGET_ROOT=`cd "$(dirname "$0")"; pwd`
#RELEASE_NAME=release_$V-$Version
RELEASE_NAME=release_$V
RELEASE_ROOT=${TARGET_ROOT}/${RELEASE_NAME}

echo RELEASE_NAME=$RELEASE_NAME

cd ${TARGET_ROOT}/
#update web dir
./replaceWebApp.sh

echo "build x86_64 updateBsp... "
make 
echo "build arm64 updateBsp... "
make updateBsp-ARM64

if [ ! -d ${RELEASE_ROOT} ]; then
	rm -rf ${RELEASE_ROOT}		
fi

[ ! -e ${RELEASE_ROOT} ] && mkdir -p ${RELEASE_ROOT}

cp -a frontend ${RELEASE_ROOT}/
cp -a conf ${RELEASE_ROOT}/
cp -a linuxPackage ${RELEASE_ROOT}/
cp  -a updateBsp updateBsp-ARM64 ${RELEASE_ROOT}/
[ -a "startup.sh" ] && cp -a startup.sh ${RELEASE_ROOT}/
[ -a "updateBsp.service" ] && cp -a updateBsp.service  ${RELEASE_ROOT}/

echo "Syncing...."
sync;sync;
sync
sync

#tar cJvf ${RELEASE_NAME}.tar.xz ${RELEASE_NAME}
tar zcvf ${RELEASE_NAME}.tar.gz ${RELEASE_NAME}

rm -rf ${RELEASE_ROOT}
echo "[Done]"

