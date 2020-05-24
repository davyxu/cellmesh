#!/bin/bash

# 默认设置代理, 国内专用
export GOPROXY=https://goproxy.io

go build -v -o ./tabtoy github.com/davyxu/tabtoy

./tabtoy -mode=v3 \
-index=Index.csv \
-go_out=../../fx/zonecfg/tab_gen.go \
-json_out=../../cfg/ZoneConfig.json \
-package=zonecfg

rm -f tabtoy