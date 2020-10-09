![.github/workflows/release.yaml](https://github.com/Faagerholm/clockify-cli/workflows/.github/workflows/release.yaml/badge.svg?branch=v1.1&event=release)	
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/4331/badge)](https://bestpractices.coreinfrastructure.org/projects/4331)	

<a href="https://www.buymeacoffee.com/Faagerholm" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

# clockify-cli
Integrate your clocking with your favorite CLI. 

### Install:

```bash
wget https://raw.githubusercontent.com/Faagerholm/clockify-cli/master/install.sh && ./install.sh
```

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
