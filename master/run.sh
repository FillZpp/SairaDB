#!/bin/sh

CURDIR=`pwd`
echo $CURDIR
OLDGOPATH=$GOPATH
export GOPATH=$CURDIR

PREFIX="/usr/local"

echo "package config
var Prefix = \"$PREFIX\"
" > src/config/prefix.go

go install master

export GOPATH=$OLDGOPATH

echo -e "\033[32mRun:\033[0m"

./bin/master $*

