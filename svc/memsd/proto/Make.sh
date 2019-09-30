#!/usr/bin/env bash

BinaryDir=../../../bin
# 协议生成
ProtoPlusGen=${BinaryDir}/protoplus
go build -v -o ${ProtoPlusGen} github.com/davyxu/protoplus

echo "生成服务器协议的go消息..."
${ProtoPlusGen} -package=sdproto -go_out=msgsvc_gen.go -genreg sd.proto