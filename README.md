![.github/workflows/release.yaml](https://github.com/Faagerholm/clockify-cli/workflows/.github/workflows/release.yaml/badge.svg?branch=v1.1&event=release)	
[![CLI Best Practices](https://bestpractices.coreinfrastructure.org/projects/4331/badge)](https://bestpractices.coreinfrastructure.org/projects/4331)	

<a href="https://www.buymeacoffee.com/Faagerholm" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>

# clockify-cli
**Version 2.0 is out!**  

Now with to new UI for a better experiance  
Integrate your clocking with your favorite CLI. 

*Please note: This tool does not take any responsibilities of spam on behalf of the user against the clockify API.*

## Install:

```bash
wget https://raw.githubusercontent.com/Faagerholm/clockify-cli/master/install.sh && ./install.sh
```

## Usage:
```
  clockify-cli [flags]  
  clockify-cli [command]
  clockift-cli // Start menu
```
## Available commands:
```
Available Commands:
  add-key         Add users API-KEY, this will store it in a yaml file.
  add-part-time   Add part-time work to your account
  check-balance   Check balance
  current-user    get current user
  default-project Select default workspace project
  help            Help about any command
  list-projects   List all projects
  menu            Select action to perform
  reset           Resets viper values
  setup           Setup
  start-timer     Select a project and start a timer
  stop-timer      Stop timer

Flags:
      --config string   config file (default is $HOME/.clockify-cli/config.yaml)
  -h, --help            help for clockify-cli
      --viper           use Viper for configuration (default true)

Use "clockify-cli [command] --help" for more information about a command.
```

## Contributing:

Please open an issue if there is something that is not working or you would like to be added to this project.

## External API:

Clockify has an API that this project heavily depends on. The API can be accessed by any user that has generated an API key from their user settings page.
More information about the API can be found here: https://clockify.me/developers-api
