#!/bin/bash
go build
VERSION=`cat bottalk-cli.go | grep 'const Version = ' | awk '{print $4}'`
echo "Publishing version ${VERSION}"
scp ./bottalk-cli frontend.bt:/home/deploy/console/console/binary.${VERSION}.bin
echo "{\"version\":${VERSION}}" > version.json
scp ./version.json frontend.bt:/home/deploy/console/console/version.json
rm -rf ./version.json
