#!/usr/bin/env bash
CURRDIR=`pwd`
cd ../../../../../..
export GOPATH=`pwd`

MsgPack=${GOPATH}/bin/msgp
go build -v -o ${MsgPack} github.com/tinylib/msgp
cd ${CURRDIR}

${MsgPack} -tests=false -io=false -file=./model/player.go