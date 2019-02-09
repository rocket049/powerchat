#!/usr/bin/env bash
mkdir -p powerchatserver/pems
cp ../server/server powerchatserver/powerchatserver
cp ../server/config.json powerchatserver/config.json
cp ../server/pems/a-cert.pem powerchatserver/pems/
cp ../server/pems/a-key.pem powerchatserver/pems/
tar cvfz powerchatserver.tar.gz powerchatserver

docker build -t powerchatserver:v$1 .
 
rm -rf powerchatserver
rm powerchatserver.tar.gz
