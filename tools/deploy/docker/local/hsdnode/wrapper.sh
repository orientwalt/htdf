#!/usr/bin/env sh

##
## Input parameters
##
## default: -x
BINARY=/root/${BINARY:-hsd}
ID=${ID:-0}
LOG=${LOG:-hsd.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'hsd' E.g.: -e BINARY=hsd_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
# echo `pwd`
# echo `ls -la ${BINARY}`
# echo `which hsd`

export HSDHOME="/root/node${ID}/.hsd"

if [ -d "`dirname ${HSDHOME}/${LOG}`" ]; then
  "$BINARY" --home "$HSDHOME" "$@" | tee "${HSDHOME}/${LOG}"
else
  "$BINARY" --home "$HSDHOME" "$@"
fi

chmod 0777 -R /root
# chown root:root -R /root