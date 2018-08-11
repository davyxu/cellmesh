#!/usr/bin/env bash
CURRDIR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`
CellMeshProtoGen=${GOPATH}/bin/cmprotogen
go build -v -o ${CellMeshProtoGen} github.com/davyxu/cellmesh/protogen
cd ${CURRDIR}
${CellMeshProtoGen} -package=proto -cmgo_out=proto_gen.go demo.proto