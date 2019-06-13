#!/usr/bin/env bash

set -e

BinaryDir=../bin

ProtocBinary=${BinaryDir}/protoc

Platform=$(go env GOHOSTOS)

if [[ ${Platform} == "linux" ]]; then
    DownloadURL=https://github.com/protocolbuffers/protobuf/releases/download/v3.8.0/protoc-3.8.0-linux-x86_64.zip
elif [[ ${Platform} == "darwin" ]]; then
    DownloadURL=https://github.com/protocolbuffers/protobuf/releases/download/v3.8.0/protoc-3.8.0-osx-x86_64.zip
elif [[ ${Platform} == "windows" ]]; then
    DownloadURL=https://github.com/protocolbuffers/protobuf/releases/download/v3.8.0/protoc-3.8.0-win64.zip
    ProtocBinary=${ProtocBinary}.exe
fi

# https://github.com/golang/protobuf/blob/master/.travis.yml
if [[ ! -f ${ProtocBinary} ]]; then
    echo "Google protocol buffer compiler is not installed, download from ${DownloadURL} and place at ${BinaryDir}"
    exit 1
fi

# cellmesh服务绑定
CellMeshProtoGen=${BinaryDir}/cmprotogen
go build -v -o ${CellMeshProtoGen} github.com/davyxu/cellmesh/tools/protogen

# 协议生成
ProtoPlusGen=${BinaryDir}/protoplus
go build -v -o ${ProtoPlusGen} github.com/davyxu/protoplus

# pb插件
GoGoFaster=${BinaryDir}/protoc-gen-gogofaster
go build -v -o ${GoGoFaster} github.com/gogo/protobuf/protoc-gen-gogofaster

# 路由工具
RouteGen=${BinaryDir}/routegen
go build -v -o ${RouteGen} github.com/davyxu/cellmesh/tools/routegen

echo "生成服务器协议的go消息..."
${ProtoPlusGen} -package=proto -go_out=msgsvc_gen.go `source ./protolist.sh svc`

echo "生成服务器协议的消息绑定..."
${CellMeshProtoGen} -package=proto -cmgo_out=msgbind_gen.go `source ./protolist.sh all`

echo "生成客户端协议的protobuf proto文件..."
${ProtoPlusGen} --package=proto -pb_out=clientmsg_gen.proto `source ./protolist.sh client`

echo "生成客户端协议的protobuf的go消息...."
${ProtocBinary} --plugin=protoc-gen-gogofaster=${GoGoFaster} --gogofaster_out=. --proto_path=. clientmsg_gen.proto

# 不使用protobuf协议文件,只使用生成后的go文件,删除之
rm -f ./clientmsg_gen.proto


echo "更新agent路由表"
${RouteGen} -package=proto -configpath=config_demo/route_rule `source ./protolist.sh client`
