#!/usr/bin/env bash

go build -v -o=${GOPATH}/bin/protoplus github.com/davyxu/protoplus


${GOPATH}/bin/protoplus -go_out=msg_gen.go -genreg -package=proto \
chat.proto \
hub.proto \
login.proto \
result.proto \
router.proto \
svc.proto