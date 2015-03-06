#!/bin/sh

CURDIR=`pwd`
echo $CURDIR
OLDGOPATH=$GOPATH
export GOPATH=$CURDIR

PREFIX="/usr/local/etc/sairadb/"

echo "package config
var prefix = \"$PREFIX\"
" > src/config/prefix.go

go install master

export GOPATH=$OLDGOPATH

echo "\033[32mRun:\033[0m"

./bin/master $*

