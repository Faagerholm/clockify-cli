package common

import (
	"fmt"
	"log"

	"github.com/Faagerholm/clockify-cli/cmd/clockify"
	"github.com/Faagerholm/clockify-cli/cmd/jira"
	"github.com/Faagerholm/clockify-cli/utils"
	"github.com/spf13/cobra"
)

var CompareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare different systems",
	Run: func(cmd *cobra.Command, args []string) {
		firstOfMonth, lastOfMonth := utils.GetFirstAndLastDayOfMonth()
		worklog, err := jira.GetWorklog(firstOfMonth, lastOfMonth, []string{"ISTEHALL", "IOTAYS"})
		if err != nil {
			log.Fatal(err)
		}
		clockifylog, _ := clockify.GetWorklog(firstOfMonth, lastOfMonth)

		c_projects := make(map[string]int)

		for _, e := range clockifylog.Entries {
			for _, c := range e.Children {
				for _, p := range c.Children {

					_, ok := c_projects[p.Project]
					if ok {
						c_projects[p.Project] += p.Duration
					} else {
						c_projects[p.Project] = p.Duration
					}
				}
			}
		}

		j_projects := make(map[string]int)

		for _, t := range worklog {
			_, ok := j_projects[t.Issue.ProjectKey]
			if ok {
				j_projects[t.Issue.ProjectKey] += t.BillableSeconds
			} else {
				j_projects[t.Issue.ProjectKey] = t.BillableSeconds
			}
		}

		fmt.Printf("From period %s to %s\n", firstOfMonth, lastOfMonth)
		fmt.Println("=============== INTERNAL ===============")
		fmt.Printf("Summary from Clockify\n\n")
		utils.DisplayProjects(c_projects)
		fmt.Println("=============== EXTERNAL ===============")
		if len(j_projects) > 0 {
			fmt.Printf("Summary from Jira\n\n")
			utils.DisplayProjects(j_projects)
		} else {
			fmt.Println("Could not load Jira projects, please check your project list and try again.")
		}

	},
}
