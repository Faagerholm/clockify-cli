#!/bin/bash

DEFAULT=$(tput sgr0)
RED=$(tput setaf 1)
GREEN=$(tput setaf 2)
BLUE=$(tput setaf 4)
LIME_YELLOW=$(tput setaf 190)

PROJECT_HOME="$HOME/.clockify-cli"

echo "${LIME_YELLOW}Initializing clockify-cli..${DEFAULT}"

if [ -d $PROJECT_HOME ]; then
        echo "Project directory already exists"
else
        $(mkdir -p $PROJECT_HOME)
        echo "$PROJECT_HOME directory is created"
fi

if [ -f "$PROJECT_HOME/config.yaml" ]; then
        echo "Config file already exists"
else
        $(touch $PROJECT_HOME/config.yaml)
        echo "config file initialized"
fi

echo "Creating alias for clockify"
if [ -f "$HOME/.zshrc" ]; then
        echo "Adding it to your zshrc file."
        alias clockify >/dev/null 2>&1 && echo "${GREEN}clockify${DEFAULT} is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >>$HOME/.zshrc
elif [ -f "$HOME/.bash_profile" ]; then
        echo "Adding it to your bash_profile file."
        alias clockify >/dev/null 2>&1 && echo "${GREEN}clockify${DEFAULT} is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >>$HOME/.bash_profile
elif [ -f "$HOME/.bashrc" ]; then
        echo "Adding it to your bashrc file."
        alias clockify >/dev/null 2>&1 && echo "${GREEN}clockify${DEFAULT} is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >>$HOME/.bashrc
elif [ -f "$HOME/.profile" ]; then
        echo "Adding it to your profile file."
        alias clockify >/dev/null 2>&1 && echo "${GREEN}clockify${DEFAULT} is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >>$HOME/.profile
elif [ -f "/etc/profile" ]; then
        # This is for the case where the user is using the system profile file
        # instead of their own.
        echo "Adding it to your default profile file."
        alias clockify >/dev/null 2>&1 && echo "${GREEN}clockify${DEFAULT} is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >>/etc/profile
else
        echo "Could not find a terminal profile, please manually add ${GREEN}alias clockify='$PROJECT_HOME/clockify-cli'${DEFAULT} to your profile."
fi

# init alias temporarily
alias clockify="$PROJECT_HOME/clockify-cli"

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

echo "----------------------------"
mv clockify-cli $PROJECT_HOME/
echo "To get started you will need a API-key. The key can be genereted on your profile page."
echo "Generate your API-key here: ${BLUE}https://clockify.me/user/settings${DEFAULT}"

# Run setup
$PROJECT_HOME/clockify-cli setup

echo "Initialization completed, please run ${GREEN}clockify help${DEFAULT} to get started."
