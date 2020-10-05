#!/bin/bash

PROJECT_HOME="$HOME/.clockify-cli"

echo 'Initializing clockify-cli'

if [ -d $PROJECT_HOME ]; then
echo "Directory already exists" ;
else
`mkdir -p $PROJECT_HOME`;
echo "$PROJECT_HOME directory is created"
fi

if [ -f "$PROJECT_HOME/config.yaml" ]; then
echo "Config file already exists" ;
else
`touch $PROJECT_HOME/config.yaml`;
echo "config file initialized"
fi

echo "Creating alias for clockify"
if [ -f "$HOME/.zshrc" ]; then
echo "Adding it to your zshrc file."
alias clockify >/dev/null 2>&1 && echo "clockify is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >> $HOME/.zshrc
elif [ -f "$HOME/.bash_profile" ]; then
echo "Adding it to your bash_profile file."
alias clockify >/dev/null 2>&1 && echo "clockify is set as an alias, skipping update of source file." || echo "alias clockify='$PROJECT_HOME/clockify-cli'" >> $HOME/.bash_profile
fi

echo "Initialization completed, please run 'clockify add-key' to get started."
echo "You will have to restart any active terminal instance to access you newly created command."
echo "You can also read the source file with e.g. 'source ~/.zshrc'".