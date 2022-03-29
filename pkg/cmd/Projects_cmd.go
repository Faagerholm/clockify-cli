package cmd

import (
	"fmt"
	"log"
	"strings"

	api "github.com/Faagerholm/clockify-cli/pkg/API"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var DefaultProjectCmd = &cobra.Command{
	Use:   "default-project",
	Short: "Select default workspace project",
	Long: `Display all workspace projects and 
	select the default project to use when starting a timer`,
	Run: func(cmd *cobra.Command, args []string) {

		CheckConfigAndPromptSetup()
		DefaultProject()
	},
}

var ListProjectsCmd = &cobra.Command{
	Use:   "list-projects",
	Short: "List all projects",
	Long:  `Display all projects in the current workspace`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckConfigAndPromptSetup()
		ListProjects()
	},
}

func ListProjects() {
	fmt.Println("Listing projects, this may take a while...")
	fmt.Println("Note, noting happens if you select a project as of now")
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
		Label:     "This is all projects I could find",
		Items:     projects,
		Templates: templates,
		Size:      20,
		Searcher:  searcher,
	}

	_, _, err = prompt.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func DefaultProject() {

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
	viper.Set("default-project", projects[i])
	viper.WriteConfig()
	fmt.Println("Default project set:", projects[i].Name)

}
