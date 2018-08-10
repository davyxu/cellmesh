#!/usr/bin/env bash
CURRDIR=`pwd`
cd ../../../../..
export GOPATH=`pwd`
if [ ! -f "${GOPATH}/bin/consul" ]; then
	go install -v github.com/hashicorp/consul
fi

cd ${GOPATH}/bin
./consul agent -server -bind=127.0.0.1 -bootstrap-expect=1  -client=127.0.0.1 -ui -data-dir=.