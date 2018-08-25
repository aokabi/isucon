#!/bin/bash

USERNAME=root
USER_HOME=/root
NODE_VERSION=v6.14.4

export PATH=$PATH:$USER_HOME/.nvm/$NODE_VERSION/bin
export NODE_PATH=$NODE_PATH:$USER_HOME/node_modules/
