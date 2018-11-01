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

# 编译
go build -v -o=${GOPATH}/src/github.com/davyxu/cellmesh/demo/bin/${Name} github.com/davyxu/cellmesh/demo/svc/${Name}

mkdir -p ${GOPATH}/src/github.com/davyxu/cellmesh/demo/bin

# 启动
${GOPATH}/src/github.com/davyxu/cellmesh/demo/bin/${Name} -logcolor