#!/usr/bin/env bash

go build -v -o=${GOPATH}/bin/protoplus github.com/davyxu/protoplus


${GOPATH}/bin/protoplus -go_out=msg_gen.go -genreg -package=meshproto msg.proto