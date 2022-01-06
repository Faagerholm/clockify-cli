package clockify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var StartActivityCmd = &cobra.Command{
	Use:   "start [project_id]",
	Short: "start timer for project. Use 'default' flag to use default project id.",
	Long: `Start timer for project. User 'default' flag to use default project id.
	You can set your default project with clockify-cli projects.
	If the flag and project id is omitted, and the default is set. Thw default will be used!`,
	Run: func(cmd *cobra.Command, args []string) {

		loc, _ := time.LoadLocation("UTC")
		cur_time := time.Now().In(loc)

		start_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", cur_time.Year(), cur_time.Month(), cur_time.Day(),
			cur_time.Hour(), cur_time.Minute(), cur_time.Second())

		project := ""
		if len(args) == 0 && viper.IsSet("default-project") {
			project = viper.GetString("default-project")
		} else if len(args) > 0 && len(args[0]) == 23 {
			project = args[0]
		} else if len(args) > 0 && (args[0] == "d" || args[0] == "default") {
			project = viper.GetString("default-project")
		} else {
			log.Fatal("Could not parse arguments. Check 'start --help' for more information.")
		}

		reqBody, err := json.Marshal(LogEntry{Start: start_time, ProjectId: project})
		if AddEntry(reqBody) != nil {
			log.Fatal(err)
		}
	},
}

var StopActivityCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops an active timer.",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.GetString("API-KEY")
		workspace := viper.GetString("WORKSPACE")
		user := viper.GetString("USER-ID")

		loc, _ := time.LoadLocation("UTC")
		cur_time := time.Now().In(loc)

		end_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", cur_time.Year(), cur_time.Month(), cur_time.Day(),
			cur_time.Hour(), cur_time.Minute(), cur_time.Second())

		reqBody, err := json.Marshal(map[string]string{
			"end": end_time,
		})

		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("https://test.clockify.me/api/v1/workspaces/%s/user/%s/time-entries", workspace, user), bytes.NewBuffer(reqBody))
		req.Header.Set("X-API-KEY", key)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Project stopped")
		}
		defer resp.Body.Close()

	},
}

var AddActivityCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an entry",
	Run: func(cmd *cobra.Command, args []string) {
		project, err := SelectProject()

		if err != nil {
			log.Fatal(err)
		}
		_ = project
	},
}

func AddEntry(body []byte) error {
	key := viper.GetString("API-KEY")
	workspace := viper.GetString("WORKSPACE")

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/time-entries", workspace), bytes.NewBuffer(body))
	req.Header.Set("X-API-KEY", key)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
