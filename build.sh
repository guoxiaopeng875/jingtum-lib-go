#!/usr/bin/env bash
if [ ! -f build.sh ]; then
    echo 'install must be run within its container folder' 1>&2
    exit 1
fi
if [ ! -d bin ]; then
      mkdir bin
fi
if [ ! -d log ]; then
      mkdir log
fi
CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"
gofmt -w src
go install testLib
export GOPATH="$OLDGOPATH"
echo 'finished'
