package cmd

import (
	"fmt"

	"github.com/Faagerholm/clockify-cli/utils"
	"github.com/spf13/cobra"
)

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update clockify-cli",
	Long:  `Update clockify-cli to the latest version.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.Update()
		if err != nil {
			fmt.Printf("Error updating clockify-cli: %s", err)
		}
	},
}
