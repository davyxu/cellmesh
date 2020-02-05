#!/usr/bin/env bash

set -e

# 生成协议
echo "Compile protoplus..."
ProtoPlus=../bin/protoplus
go build -v -o=${ProtoPlus} github.com/davyxu/protoplus

echo "Gen proto..."
${ProtoPlus} -go_out=msg_gen.go -genreg -package=proto `source ./protolist.sh all`


# 生成网关路由
echo "Compile routegen..."
RouteGen=../bin/routegen
go build -v -o=${RouteGen} github.com/davyxu/cellmesh/tool/routegen

echo "Upload route config to discovery..."
${RouteGen} `source ./protolist.sh client`