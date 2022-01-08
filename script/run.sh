#!/bin/bash

SCRIPT=`realpath $0`
SCRIPTPATH=`dirname $SCRIPT`
PROJECTDIR="$(dirname "$SCRIPTPATH")"

cd ${PROJECTDIR}/cmd/wasm
echo "Building wasm"
GOOS=js GOARCH=wasm go build -o ${PROJECTDIR}/web/trace.wasm

cd ${PROJECTDIR}/cmd/server
echo "Building server"
go build
./server 8080