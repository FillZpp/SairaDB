#!/bin/sh

CURDIR=`pwd`
echo $CURDIR
OLDGOPATH=$GOPATH
export GOPATH=$CURDIR

go install master

export GOPATH=$OLDGOPATH

echo "\033[32mRun:\033[0m"

./bin/master

