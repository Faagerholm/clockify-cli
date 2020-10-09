#!/bin/bash

PROJECT_HOME="$HOME/.clockify-cli"

echo 'Initializing clockify-cli'

if [ -d $PROJECT_HOME ]; then
        echo "Directory already exists"
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
        # alias clockify >/dev/null 2>&1 && echo "clockify is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >> $HOME/.zshrc
elif [ -f "$HOME/.bash_profile" ]; then
        echo "Adding it to your bash_profile file."
        alias clockify >/dev/null 2>&1 && echo "clockify is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >>$HOME/.bash_profile
fi

os="$OSTYPE"
processor="$(uname -m)"
echo 'Downloading the latest release for your machine...'
# Download the right binary, based on your OS.
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        if [[ "$processor" == "i386" ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'linux-386' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        elif [[ "$processor" == *"x86"* ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'linux-amd64' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        fi
elif [[ "$OSTYPE" == "darwin"* ]]; then
        if [[ "$processor" == "i386" ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'darwin-386' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        elif [[ "$processor" == *"x86"* ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'darwin-amd64' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        fi
elif [[ "$OSTYPE" == "cygwin" ]]; then
        # POSIX compatibility layer and Linux environment emulation for Windows
        if [[ "$processor" == "i386" ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'windows-386' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        elif [[ "$processor" == *"x86"* ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'windows-amd64' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        fi
elif [[ "$OSTYPE" == "msys" ]]; then
        # Lightweight shell and GNU utilities compiled for Windows (part of MinGW)
        if [[ "$processor" == "i386" ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'windows-386' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        elif [[ "$processor" == *"x86"* ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'windows-amd64' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        fi
elif [[ "$OSTYPE" == "win32" ]]; then
        # I'm not sure this can happen.
        if [[ "$processor" == "i386" ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'windows-386' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        elif [[ "$processor" == *"x86"* ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'windows-amd64' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        fi
elif [[ "$OSTYPE" == "freebsd"* ]]; then
        # ...
        if [[ "$processor" == "i386" ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'linux-386' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        elif [[ "$processor" == *"x86"* ]]; then
                curl -s https://api.github.com/repos/faagerholm/clockify-cli/releases/latest |
                        grep browser_download_url |
                        cut -d '"' -f 4 |
                        grep 'linux-amd64' |
                        wget -qi-
                tarfilename="$(find . -name "*.tar.gz")"
                tar -xzf $tarfilename && mv $tarfilename/clockify-cli clockify-cli
                sudo rm $tarfilename
        fi
else
        # Unknown.
        echo 'unknows operating system, unable to download binary.\nYou can visit https://github.com/Faagerholm/clockify-cli/releases/latest to download the right binary'
        echo 'Please make a issue explaining the problem and I will fix it for you. This should not happen!'
fi

echo 'Done!'

mv clockify-cli $PROJECT_HOME/clockify-cli

$PROJECT_HOME/clockify-cli init

echo "Initialization completed, please run 'clockify add-key' to get started."
echo "You will have to restart any active terminal instance to access you newly created command."
echo "You can also read the source file with e.g. 'source ~/.zshrc'".

