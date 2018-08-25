#!/bin/bash

APPDIR=$(dirname $0)/..

if [ -f $APPDIR/../../standalone/env.sh ]; then
    . $APPDIR/../../standalone/env.sh
else
    export PATH=$PATH:/root/.nvm/v6.14.4/bin
    export NODE_PATH=/root/node_modules/
fi

cd $APPDIR
exec node master.js
