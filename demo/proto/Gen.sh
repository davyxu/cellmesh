#!/usr/bin/env bash

CURRDIR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`

set -e

CellMeshProtoGen=${GOPATH}/bin/cmprotogen
go build -v -o ${CellMeshProtoGen} github.com/davyxu/cellmesh/tools/protogen

RouteGen=${GOPATH}/bin/routegen
go build -v -o ${RouteGen} github.com/davyxu/cellmesh/tools/routegen

cd ${CURRDIR}

${CellMeshProtoGen} -package=proto -cmgo_out=proto_gen.go `source ./protolist.sh`

${RouteGen} `source ./protolist.sh`