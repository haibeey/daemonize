#!/bin/bash

OS_TYPE=$(uname | tr '[:upper:]' '[:lower:]')

wget -O daemonize.gz https://github.com/haibeey/daemonize/releases/download/initial/daemonize.$OS_TYPE.gz

gunzip -dc daemonize.gz > daemonize && chmod +x daemonize

if [ "$OS_TYPE" == "linux" ]; then
    sudo mv daemonize /usr/bin/
elif [ "$OS_TYPE" == "darwin" ]; then
    sudo mv daemonize /usr/local/bin/
else
    echo "un supported  os: $OS_TYPE"
fi

rm daemonize.gz