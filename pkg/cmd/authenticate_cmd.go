package cmd

import (
	"fmt"
	"log"
	"strings"

	api "github.com/Faagerholm/clockify-cli/pkg/API"
	utils "github.com/Faagerholm/clockify-cli/pkg/Utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var GetUserCmd = &cobra.Command{
	Use:   "current-user",
	Short: "get current user",
	Run: func(cmd *cobra.Command, args []string) {
		getUser()
	},
}

var AddKeyCmd = &cobra.Command{
	Use:   "add-key [API-KEY]",
	Short: "Add users API-KEY, this will store it in a yaml file.",
	Long: `Add users API-KEY, get the key from clockify.me/user/settings.
	At the bottom of the page, generate KEY.`,
	Run: func(cmd *cobra.Command, args []string) {
		key := ""
		fmt.Println(len(args))
		if len(args) == 0 {
			key = viper.GetString("API-KEY")
		} else {
			key = args[0]
		}
		AddKey(key)
	},
}

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup",
	Long:  `Setup the application, this will ask you for your API key and store it in a yaml file.`,
	Run: func(cmd *cobra.Command, args []string) {
		Authenticate()
	},
}

func Authenticate() {

	// This is a method for setting the API key and storing it in a yaml file
	key := viper.GetString("API-KEY")

	prompt := promptui.Prompt{
		Label:     "Do you wish to add a new API key?",
		IsConfirm: true,
	}

	result, err := prompt.Run()

	if err != nil {
		// fmt.Printf("Prompt failed %v\n", err)
		return
	}

	if strings.ToLower(result) == "y" {
		prompt := promptui.Prompt{
			Label:     "Enter your API key",
			Default:   key,
			AllowEdit: true,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		viper.Set("API-KEY", result)
		viper.WriteConfig()

		user := api.GetUser()
		if user != nil {
			viper.Set("USER-ID", user.ID)
			viper.Set("WORKSPACE", user.ActiveWorkspace)

			fmt.Println("Updating config with user id.")
			fmt.Println("You're ready to go.. check help for more commands")
			viper.WriteConfig()
		} else {
			log.Println("Could not find user.. try again later.")
		}
	}

}

func getUser() {
	user := api.GetUser()
	if user != nil {
		workspace := utils.ExtractWorksapce()
		if user.ActiveWorkspace != workspace {
			viper.Set("WORKSPACE", user.ActiveWorkspace)
			viper.WatchConfig()
		}
		fmt.Printf("Active user: %s\n", user.Name)
	} else {
		fmt.Println("Could not find user.. try again later.")
	}
}

func AddKey(key string) {
	prompt := promptui.Prompt{
		Label:     "Enter your API key",
		Default:   key,
		AllowEdit: true,
	}

	result, err := prompt.Run()

	if err != nil {
		// fmt.Printf("Prompt failed %v\n", err)
		return
	}
	key = result

	prompt = promptui.Prompt{
		Label:     "Do you wish to add a new API key",
		IsConfirm: true,
	}

	result, err = prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch strings.ToLower(result) {
	case "y":
		viper.Set("API-KEY", key)
		fmt.Println("Saving", viper.Get("API-KEY"), `as your user key, this can be changed later by initializing the same command. 
As of now, no more the one key can be used at the same time.`)
		err := viper.WriteConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error in config file: %s \n", err))
		}

	case "n":
		fmt.Println("The key was NOT added.")
	}
}
