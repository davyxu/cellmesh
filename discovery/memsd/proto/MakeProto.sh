#!/usr/bin/env bash

CURRDIR=`pwd`
cd ../../../../../../..
export GOPATH=`pwd`

set -e
Protoc=${GOPATH}/bin/protoc

# cellmesh服务绑定
CellMeshProtoGen=${GOPATH}/bin/cmprotogen
go build -v -o ${CellMeshProtoGen} github.com/davyxu/cellmesh/tools/protogen

# 协议生成
ProtoPlusGen=${GOPATH}/bin/protoplus
go build -v -o ${ProtoPlusGen} github.com/davyxu/protoplus

cd ${CURRDIR}

# windows下时，添加后缀名
if [ `go env GOHOSTOS` == "windows" ];then
	EXESUFFIX=.exe
fi

echo "生成服务器协议的go消息..."
${ProtoPlusGen} -package=proto -go_out=msgsvc_gen.go `source ./protolist.sh svc`

echo "生成服务器协议的消息绑定..."
${CellMeshProtoGen} -package=proto -cmgo_out=msgbind_gen.go `source ./protolist.sh svc`