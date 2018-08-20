#!/usr/bin/env bash
Name=$1

if [ "${Name}" == "" ]
then
	echo "Usage: RunDemoSvc.sh name"
	exit 1
fi

# 错误退出
set -e

# 编译
go build -v -o=${GOPATH}/bin/${Name} github.com/davyxu/cellmesh/demo/${Name}

# 启动
${GOPATH}/bin/${Name} -colorlog