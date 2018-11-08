#!/usr/bin/env bash
Name=$1

if [ "${Name}" == "" ]
then
	echo "Usage: RunDemoSvc.sh name"
	exit 1
fi

CURRDIR=`pwd`
cd ../../../../..
export GOPATH=`pwd`

# 错误退出
set -e

mkdir -p ${GOPATH}/src/github.com/davyxu/cellmesh/demo/bin

# 工作路径需要在bin下
cd ${GOPATH}/src/github.com/davyxu/cellmesh/demo/bin

# 编译
go build -v -o=./${Name} github.com/davyxu/cellmesh/demo/svc/${Name}

# 在bin下启动服务器
./${Name}