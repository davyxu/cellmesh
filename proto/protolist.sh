#!/usr/bin/env bash

# protolist.sh all 将proto_client.txt和proto_svc.txt的内容输出为行
# protolist.sh XX  将proto_XX.txt的内容输出为行

if [ "$1" == "all" ]
then
{ grep -o '^[^#]*' proto_client.txt; echo " ";grep -o '^[^#]*' proto_svc.txt; }| tr -s "\r\n" " "
else
{ grep -o '^[^#]*' proto_${1}.txt; echo " "; }| tr -s "\r\n" " "
fi