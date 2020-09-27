#!/bin/bash

PROJECT_HOME="/Users/$USER/.test-clockify-cli"

echo 'Initializing clockify-cli'

if [ -d $PROJECT_HOME ]; then
echo "Directory already exists" ;
else
`mkdir -p $PROJECT_HOME`;
echo "$PROJECT_HOME directory is created"
fi

if [ -f "$PROJECT_HOME/config.json" ]; then
echo "Config file already exists" ;
else
`touch $PROJECT_HOME/config.json`;
echo "config file initialized"
fi
