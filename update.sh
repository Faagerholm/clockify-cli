#!/bin/bash

DEFAULT=$(tput sgr0)
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
BLUE=$(tput setaf 4)
LIME_YELLOW=$(tput setaf 190)

PROJECT_HOME="$HOME/.clockify-cli"

echo "${LIME_YELLOW}Initializing clockify-cli..${DEFAULT}"

if [ ! -d $PROJECT_HOME ]; then
    echo "Project directory does not exist, please run install.sh first"
fi

VERSION=$(curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest | grep tag_name | cut -d '"' -f 4)
currentVersion=$(cat $PROJECT_HOME/.version)

if [ "$VERSION" == "$currentVersion" ]; then
    echo "You are already on the latest version"
    exit 0
fi

echo "Downloading version $VERSION"

PROCESSOR="$(uname -m)"
OS_PROCESSOR=""
echo "${LIME_YELLOW}Downloading the latest release for your machine..${DEFAULT}"

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    if [[ "$PROCESSOR" == "i386" ]]; then
        OS_PROCESSOR="linux-386"
    elif [[ "$PROCESSOR" == *"x86"* ]]; then
        OS_PROCESSOR="linux-amd64"
    fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
    if [[ "$PROCESSOR" == "i386" ]]; then
        OS_PROCESSOR="darwin-386"
    elif [[ "$PROCESSOR" == *"x86"* ]]; then
        OS_PROCESSOR="darwin-amd64"
    fi
elif [[ "$OSTYPE" == "cygwin" ]]; then
    # POSIX compatibility layer and Linux environment emulation for Windows
    if [[ "$PROCESSOR" == "i386" ]]; then
        OS_PROCESSOR="windows-386"
    elif [[ "$PROCESSOR" == *"x86"* ]]; then
        $OS_PROCESSOR="windows-amd64"
    fi
elif [[ "$OSTYPE" == "msys" ]]; then
    # Lightweight shell and GNU utilities compiled for Windows (part of MinGW)
    if [[ "$PROCESSOR" == "i386" ]]; then
        OS_PROCESSOR="windows-386"
    elif [[ "$PROCESSOR" == *"x86"* ]]; then
        OS_PROCESSOR="windows-amd64"
    fi
elif [[ "$OSTYPE" == "win32" ]]; then
    # I'm not sure this can happen.
    if [[ "$PROCESSOR" == "i386" ]]; then
        OS_PROCESSOR="windows-386"
    elif [[ "$PROCESSOR" == *"x86"* ]]; then
        OS_PROCESSOR="windows-amd64"
    fi
elif [[ "$OSTYPE" == "freebsd"* ]]; then
    # ...
    if [[ "$PROCESSOR" == "i386" ]]; then
        OS_PROCESSOR="linux-386"
    elif [[ "$PROCESSOR" == *"x86"* ]]; then
        OS_PROCESSOR="linux-amd64"
    fi
fi

curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
    grep browser_download_url |
    cut -d '"' -f 4 |
    grep $OS_PROCESSOR |
    wget -qi-
tarfilename="$(find . -name "*$OS_PROCESSOR*.tar.gz")"
tar -xzf $tarfilename
rm $tarfilename
rm "$tarfilename.md5"

echo "${LIME_YELLOW}Replacing old version..${DEFAULT}"
rm -rf $PROJECT_HOME/bin
mv clockify-cli $PROJECT_HOME/bin

# save git version to config directory
curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
    grep tag_name |
    cut -d '"' -f 4 |
    sed 's/v//' >$PROJECT_HOME/.version
