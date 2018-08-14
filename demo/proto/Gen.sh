#!/usr/bin/env bash

CURRDIR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`

set -e

CellMeshProtoGen=${GOPATH}/bin/cmprotogen
go build -v -o ${CellMeshProtoGen} github.com/davyxu/cellmesh/protogen

cd ${CURRDIR}

${CellMeshProtoGen} -package=proto -cmgo_out=proto_gen.go \
demo.proto \
router.proto