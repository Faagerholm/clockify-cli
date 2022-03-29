package cmd

import (
	"fmt"
	"time"

	api "github.com/Faagerholm/clockify-cli/pkg/API"
	model "github.com/Faagerholm/clockify-cli/pkg/Model"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var StopTimerCmd = &cobra.Command{
	Use:   "stop-timer",
	Short: "Stop timer",
	Long:  `Stop the current timer`,
	Run: func(cmd *cobra.Command, args []string) {

		CheckConfigAndPromptSetup()
		StopTimer()
	},
}

func StopTimer() {
	loc, _ := time.LoadLocation("UTC")
	cur_time := time.Now().In(loc)
	end_time_str := cur_time.Format("2006-01-02T15:04:05.000Z")

	entry := api.Stop(end_time_str)

	prompt := promptui.Prompt{
		Label: "Please fill in a description",
	}

	description, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	entry.Description = description
	updateEntry := convertyEntryToUpdateEntry(entry)
	api.AddDescription(entry.ID, updateEntry)
}

func convertyEntryToUpdateEntry(entry model.Entry) model.UpdateEntry {
	return model.UpdateEntry{
		Start:       entry.TimeInterval.Start,
		Billable:    entry.Billable,
		Description: entry.Description,
		ProjectID:   entry.ProjectID,
		TaskID:      entry.TaskID,
		End:         entry.TimeInterval.End,
		TagIDs:      entry.TagIDs,
	}
}
