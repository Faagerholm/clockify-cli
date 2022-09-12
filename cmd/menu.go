package cmd

import (
	"fmt"

	"github.com/Faagerholm/clockify-cli/domain"
	"github.com/Faagerholm/clockify-cli/utils"
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
		Items:     domain.MenuActions,
		Templates: templates,
		Size:      10,
	}

	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch domain.MenuActions[i] {
	case domain.MenuActionChangeAPIKey:
		Authenticate()
	case domain.MenuActionStart:
		CheckConfigAndPromptSetup()
		StartProject()
	case domain.MenuActionStop:
		CheckConfigAndPromptSetup()
		StopTimer()
	case domain.MenuActionShowProjects:
		CheckConfigAndPromptSetup()
		ListProjects()
	case domain.MenuActionCheckBalance:
		CheckConfigAndPromptSetup()
		CheckBalance()
	case domain.MenuActionSetPartTime:
		AddPartTimeTimespan()
	case domain.MenuActionVerifyMonth:
		CheckConfigAndPromptSetup()
		VerifyFullMonth()
	case domain.MenuActionUpdate:
		utils.Update()
	case domain.MenuActionQuit:
		fmt.Println(utils.RandomExitGreeting())
	default:
		fmt.Println("Unknown action")
	}
}
