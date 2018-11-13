#!/usr/bin/env bash

if [ "$1" == "all" ]
then
{ grep -o '^[^#]*' proto_client.txt; echo " ";grep -o '^[^#]*' proto_svc.txt; }| tr -s "\r\n" " "
else
{ grep -o '^[^#]*' proto_${1}.txt; echo " "; }| tr -s "\r\n" " "
fi