package cmd

import (
	"fmt"

	model "github.com/Faagerholm/clockify-cli/pkg/Model"
	utils "github.com/Faagerholm/clockify-cli/pkg/Utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var MenuCmd = &cobra.Command{
	Use:   "menu",
	Short: "Select action to perform",
	Long: `Display all available actions and
	select the action to perform`,
	Run: func(cmd *cobra.Command, args []string) {
		Menu()
	},
}

func Menu() {

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F449 {{ .Name | cyan }}",
		Inactive: "  {{ .Name | cyan }}",
		Selected: "\U0001F449 {{ .Name | red | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select action",
		Items:     model.MenuActions,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch model.MenuActions[i] {
	case model.MenuActionChangeAPIKey:
		Authenticate()
	case model.MenuActionStart:
		StartProject()
	case model.MenuActionStop:
		StopTimer()
	case model.MenuActionShowProjects:
		ListProjects()
	case model.MenuActionCheckBalance:
		CheckBalance()
	case model.MenuActionSetPartTime:
		AddPartTimeTimespan()
	case model.MenuActionVerifyMonth:
		VerifyFullMonth()
	case model.MenuActionQuit:
		fmt.Println(utils.RandomExitGreeting())
	default:
		fmt.Println("Unknown action")
	}
}
