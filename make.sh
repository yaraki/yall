#!/bin/sh
# Requires Godag for build
# http://code.google.com/p/godag/

if [ $# -lt 1 ] ; then
    gd src -o yall
elif [ $1 = 'test' ] ; then
    gd src -test
elif [ $1 = 'run' ] ; then
    gd src -o yall && rlwrap ./yall
elif [ $1 = 'fmt'] ; then
    gd src -fmt
fi
