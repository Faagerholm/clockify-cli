# clockify-cli
Integrate your clocking with your favorite CLI. 

### Usage:
```
  clockify-cli [flags]  
  clockify-cli [command]
```
### Available commands:
```
  add-key      Add users API-KEY
  balance      Display if you're above or below zero balance.
  help         Help about any command
  off-projects Select which projects should be omitted from reports
  projects     Select default workspace project
  reset        Resets viper values
  start        start timer for project. Use 'default' flag to use default project id.
  stop         Stops an active timer.
  user         get current user
  version      Print the version number of Clockify-cli
  workspace    Get workspaces

Flags:
      --config string   config file (default is $HOME/.clockify-cli/config.yaml) (default "./config.yaml")
  -h, --help            help for clockify-cli
      --viper           use Viper for configuration (default true)
```