#!/usr/bin/env bash

# linq查询
go get -v github.com/ahmetb/go-linq

# proto文件解析
go get -v github.com/davyxu/protoplus/codegen

# 词法器
go get -v github.com/davyxu/golexer

# 网络库 消息绑定
go get -v github.com/davyxu/cellnet/protoc-gen-msg

# 网络库
go get -v github.com/davyxu/golog

# 二进制编解码
go get -v github.com/davyxu/goobjfmt

# gogo Protobuf 的protoc插件,用于go源码生成
go get -v github.com/gogo/protobuf/protoc-gen-gogofaster

# 预编译demo
go get -v github.com/davyxu/cellmesh/demo/svc/agent
go get -v github.com/davyxu/cellmesh/demo/svc/game
go get -v github.com/davyxu/cellmesh/demo/svc/hub
go get -v github.com/davyxu/cellmesh/demo/svc/login
go get -v github.com/davyxu/cellmesh/demo/svc/client