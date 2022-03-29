package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	api "github.com/Faagerholm/clockify-cli/pkg/API"
	model "github.com/Faagerholm/clockify-cli/pkg/Model"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var StartTimerCmd = &cobra.Command{
	Use:   "start-timer",
	Short: "Select a project and start a timer",
	Long: `Display all projects in the current workspace and
	select the project to start a timer in. A default project can be used`,
	Run: func(cmd *cobra.Command, args []string) {

		CheckConfigAndPromptSetup()
		StartProject()
	},
}

func StartProject() {

	project := checkDefaultProject()

	if project == nil {

		projects, err := api.GetProjects()
		if err != nil {
			log.Fatal(err)
		}

		templates := &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "\U0001F449 {{ .Name | cyan }}",
			Inactive: "  {{ .Name | cyan }}",
			Selected: "\U0001F449 {{ .Name | red | cyan }}",
		}

		searcher := func(input string, index int) bool {
			project := projects[index]
			name := strings.Replace(strings.ToLower(project.Name), " ", "", -1)
			input = strings.Replace(strings.ToLower(input), " ", "", -1)

			return strings.Contains(strings.ToLower(name), input)
		}
		prompt := promptui.Select{
			Label:     "Select default project (this is used when starting a timer)",
			Items:     projects,
			Templates: templates,
			Size:      10,
			Searcher:  searcher,
		}

		i, _, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		project = &projects[i]
	}
	loc, _ := time.LoadLocation("UTC")
	now := time.Now().In(loc)

	api.Start(now, project.ID)

}

func checkDefaultProject() *model.Project {
	var defaultProject *model.Project
	viper.UnmarshalKey("default-project", &defaultProject)
	if defaultProject == nil {
		fmt.Println("No default project set")
		return nil
	}

	prompt := promptui.Select{
		Label: "Do you wish to use your default project?",
		Items: []string{"Yes", "No"},
	}
	i, _, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return nil
	}
	if i == 1 {
		return nil
	}

	return defaultProject
}
