#!/usr/bin/env bash
CURRDIR=`pwd`
cd ../../../../..
export GOPATH=`pwd`

go install -v github.com/davyxu/cellmesh/discovery/memsd

cd ${GOPATH}/bin
./memsd -addr=127.0.0.1:8900