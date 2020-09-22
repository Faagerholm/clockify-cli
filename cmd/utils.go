package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Clockify-cli",
	Long:  `All software has versions. This is clockify-cli's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Clockify-cli application version 0.1 -- HEAD")
	},
}
