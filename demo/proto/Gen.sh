#!/usr/bin/env bash

CURRDIR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`

set -e

CellMeshProtoGen=${GOPATH}/bin/cmprotogen
go build -v -o ${CellMeshProtoGen} github.com/davyxu/cellmesh/tools/protogen

ProtoPlusGen=${GOPATH}/bin/protoplus
go build -v -o ${ProtoPlusGen} github.com/davyxu/protoplus

RouteGen=${GOPATH}/bin/routegen
go build -v -o ${RouteGen} github.com/davyxu/cellmesh/tools/routegen

cd ${CURRDIR}

echo "Generating proto..."
${CellMeshProtoGen} -package=proto -cmgo_out=msgbind_gen.go `source ./protolist.sh`

${ProtoPlusGen} -package=proto -go_out=proto_gen.go `source ./protolist.sh`

echo "Uploading route table..."
${RouteGen} -configpath=cm_demo/config/agent/route_rule `source ./protolist.sh`