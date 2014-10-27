#!/bin/bash
#
# Source this file to set up the correct environment variables
#

PWD=`pwd`
APPENGINE=$PWD/go_appengine
VENDOR=$PWD/.vendor

export GOROOT=$APPENGINE/goroot
export GOPATH=$VENDOR:$PWD
export PATH=$APPENGINE:$PWD/bin:$VENDOR/bin:$PATH
