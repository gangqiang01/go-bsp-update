#!/bin/bash
#v2.0
TARGET_DIR=`cd "$(dirname "$0")"; pwd`
productName="updateBsp"
## Environment
## diffrent platform has diffrent env value.
## setting for display manager's xauthority.
# exportVar="/run/user/1000/gdm/Xauthority"
# export XAUTHORITY=$exportVar

#for screenshot function
#export DISPLAY=:0
## this is used for ubuntu20.04 and later.
#export NO_MIT_SHM=1

## platform
## different platform to run different bin

cd $TARGET_DIR

key=`uname -m`
if [[ "$key" = "x86_64" ]]; then
	./${productName}
elif [[ "$key" = "arm" || "$key" = "armv7l" ]];then
	./${productName}-ARM
elif [[ "$key" = "aarch64_be"  ||  "$key" = "aarch64"  ||  "$key" = "armv8b"  ||  "$key" = "armv8l" ]];then
	./${productName}-ARM64
else
	echo "[Warning]Maybe Not Support."
fi

sleep 2

