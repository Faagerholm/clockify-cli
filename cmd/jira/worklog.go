package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var WorklogCmd = &cobra.Command{
	Use:   "jira-hours",
	Short: "Get your Jira hours.",
	Run: func(cmd *cobra.Command, args []string) {

		from, err := time.Parse("YYYY-MM-DD", "2021-11-01")
		if err != nil {
			log.Fatal(err)
			return
		}
		to, err := time.Parse("YYYY-MM-DD", "2021-11-30")
		if err != nil {
			log.Fatal(err)
			return
		}

		worklog, err := GetWorklog(from, to, []string{"ISTEHALL", "IOTAYS"})
		if err != nil {
			log.Fatal(err)
		} else {
			log.Println(worklog)
		}
	},
}

func GetWorklog(from time.Time, to time.Time, worker []string) ([]worklogEntry, error) {
	jira_url := viper.GetString("JIRA-URL")
	personal_token := viper.GetString("jira-personal-token")

	if len(jira_url) == 0 {
		return nil, errors.New("JIRA PROJECT URL NOT SET! Please check your config file")
	}
	if len(personal_token) == 0 {
		return nil, errors.New("JIRA USER TOKEN NOT SETUP PROPERLY")
	}

	body, err := json.Marshal(requestBody{from.String(), to.String(), worker})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://%s/rest/tempo-timesheets/4/worklogs/search", jira_url), bytes.NewBuffer(body))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", personal_token))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res []worklogEntry

	json.NewDecoder(resp.Body).Decode(&res)
	return res, err
}
