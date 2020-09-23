package cmd

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

var startActivityCmd = &cobra.Command{
	Use:   "start [project_id]",
	Short: "start timer for project. Use 'default' flag to use default project id.",
	Long: `Start timer for project. User 'default' flag to use default project id.
	You can set your default project with clockify-cli projects.
	If the flag and project id is omitted, and the default is set. Thw default will be used!`,
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		workspace := viper.Get("WORKSPACE")

		loc, _ := time.LoadLocation("UTC")
		cur_time := time.Now().In(loc)

		start_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", cur_time.Year(), cur_time.Month(), cur_time.Day(),
			cur_time.Hour(), cur_time.Minute(), cur_time.Second())

		project := ""
		if len(args) == 0 && viper.IsSet("default-project") {
			project = viper.Get("default-project").(string)
		} else if len(args) > 0 && len(args[0]) == 23 {
			project = args[0]
		} else if len(args) > 0 && (args[0] == "d" || args[0] == "default") {
			project = viper.Get("default-project").(string)
		} else {
			fmt.Println("Could not parse arguments. Check 'start --help' for more information.")
			return
		}

		reqBody, err := json.Marshal(map[string]string{
			"start":     start_time,
			"projectId": project,
		})

		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.clockify.me/api/v1/workspaces/%s/time-entries", workspace), bytes.NewBuffer(reqBody))
		req.Header.Set("X-API-KEY", key)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println("Project started")
		}
		defer resp.Body.Close()
	},
}

var stopActivityCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops an active timer.",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		workspace := viper.Get("WORKSPACE")
		user := viper.Get("USER-ID")

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

func getProjects() {

}
