package jira

import (
	"log"

	"github.com/spf13/cobra"
)

var ImportHoursCmd = &cobra.Command{
	Use:   "import-jira-hours",
	Short: "Comapre your Jira hours with clockify.",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("I'm not implemented yet")
	},
}
