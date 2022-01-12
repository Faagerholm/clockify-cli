package controller

import (
	"fmt"

	model "github.com/Faagerholm/clockify-cli/pkg/Model"
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
		Items:     model.MainMenuActions,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch model.MainMenuActions[i] {
	case model.MainMenuActionChangeAPIKey:
		Authenticate()
	case model.MainMenuActionStart:
		StartProject()
	case model.MainMenuActionStop:
		StopTimer()
	case model.MainMenuActionShowProjects:
		ListProjects()
	case model.MainMenuActionCheckBalance:
		CheckBalance()
	case model.MainMenuActionSetPartTime:
		AddPartTimeTimespan()
	case model.MainMenuActionQuit:
		fmt.Println("Bye!")
	default:
		fmt.Println("Unknown action")
	}
}
