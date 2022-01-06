#!/bin/bash

SCRIPT=`realpath $0`
SCRIPTPATH=`dirname $SCRIPT`
PROJECTDIR="$(dirname "$SCRIPTPATH")"

cd ${PROJECTDIR}/cmd/trace
echo "Building trace wasm"
GOOS=js GOARCH=wasm go build -o ${PROJECTDIR}/web/trace.wasm

cd ${PROJECTDIR}/cmd/server
echo "Building server"
go build
./server 8080