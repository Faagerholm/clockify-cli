![.github/workflows/release.yaml](https://github.com/Faagerholm/clockify-cli/workflows/.github/workflows/release.yaml/badge.svg?branch=v1.1&event=release)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/4331/badge)](https://bestpractices.coreinfrastructure.org/projects/4331)

# clockify-cli
Integrate your clocking with your favorite CLI. 

### Usage:
```
  clockify [flags]  
  clockify [command]
```
### Available commands:
```
Available Commands:
  add-key      Add users API-KEY
  balance      Display if you're above or below zero balance.
  help         Help about any command
  init
  off-projects Select which projects should be omitted from reports
  projects     Select default workspace project
  reset        Resets viper values
  start        start timer for project. Use 'default' flag to use default project id.
  stop         Stops an active timer.
  user         get current user
  version      Print the version number of Clockify-cli
  workspace    Get workspaces

Flags:
      --config string   config file (default is $HOME/.clockify-cli/config.yaml) (default "/Users/fagerholm/.clockify-cli/config.yaml")
  -h, --help            help for clockify-cli
      --viper           use Viper for configuration (default true)

Use "clockify-cli [command] --help" for more information about a command.
```