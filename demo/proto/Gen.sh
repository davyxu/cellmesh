#!/usr/bin/env bash
CURRDIR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`
CellMeshProtoGen=${GOPATH}/bin/cmprotogen
go build -v -o ${CellMeshProtoGen} github.com/davyxu/cellmesh/protogen
${CellMeshProtoGen} -package=proto -cmgo_out=${CURRDIR}/proto_gen.go ${CURRDIR}/demo.proto