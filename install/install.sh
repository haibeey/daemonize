#!/bin/bash

OS_TYPE=$(uname | tr '[:upper:]' '[:lower:]')

wget -O daemonize.gz https://github.com/haibeey/daemonize/releases/download/initial/daemonize.$OS_TYPE.gz

gunzip -dc daemonize.gz > daemonize && chmod +x daemonize

# Check if sudo is installed
if command -v sudo > /dev/null 2>&1; then
    USE_SUDO="sudo"
else
    USE_SUDO=""
fi

if [ "$OS_TYPE" == "linux" ]; then
    $USE_SUDO mv daemonize /usr/bin/
elif [ "$OS_TYPE" == "darwin" ]; then
    $USE_SUDO mv daemonize /usr/local/bin/
else
    echo "Unsupported OS: $OS_TYPE"
fi

rm daemonize.gz