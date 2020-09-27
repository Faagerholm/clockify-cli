package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/*
  "dateRangeStart": "2020-05-10T00:00:00.000Z",
  "dateRangeEnd": "2020-10-16T23:59:59.000Z",
  "detailedFilter": {},
  "exportType": "JSON",
  "rounding": "false",
  "amountShown": "HIDE_AMOUNT",
  "users": {
    "ids": ["5f0c5965f93abb3018f3aa10"],
    "contains": "CONTAINS",
    "status": "ALL"
  }
*/

type result struct {
	Entries []entry `json:"timeentries"`
}

type entry struct {
	ProjectId string    `json:"projectId"`
	Start     time.Time `json:"timeInterval.start"`
	End       time.Time `json:"timeInterval.end"`
	Duration  int       `json:"timeInterval.duration"`
}

type report struct {
	Start          string         `json:"dateRangeStart"`
	End            string         `json:"dateRangeEnd"`
	DetailedFilter *report_filter `json:"detailedFilter,omitempty"`
	ExportType     string         `json:"exportType"`
	Rounding       string         `json:"rounding"`
	AmountShown    string         `json:"amountShown"`
	Users          *report_user   `json:"users"`
}
type report_filter struct {
}
type report_user struct {
	Ids      []string `json:"ids"`
	Contains string   `json:"contains"`
	Status   string   `json:"status"`
}

var saldoCmd = &cobra.Command{
	Use:   "saldo",
	Short: "Display if you're above or below zero saldo.",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		workspace := viper.Get("WORKSPACE")
		user := viper.GetString("USER-ID")

		loc, _ := time.LoadLocation("UTC")
		cur_time := time.Now().In(loc)

		start_time := fmt.Sprintf("%d-01-01T00:00:00Z", cur_time.Year())
		end_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", cur_time.Year(), cur_time.Month(), cur_time.Day(),
			cur_time.Hour(), cur_time.Minute(), cur_time.Second())

		users := report_user{
			Ids:      []string{user},
			Contains: "CONTAINS",
			Status:   "All",
		}
		filter := report_filter{}
		reqBody, err := json.Marshal(report{
			Start:          start_time,
			End:            end_time,
			DetailedFilter: &filter,
			ExportType:     "JSON",
			Rounding:       "false",
			AmountShown:    "HIDE_AMOUNT",
			Users:          &users,
		})
		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		req, _ := http.NewRequest("POST", fmt.Sprintf("https://reports.api.clockify.me/v1/workspaces/%s/reports/detailed", workspace), bytes.NewBuffer(reqBody))
		req.Header.Set("X-API-KEY", key)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		} else {
			var dat map[string]interface{}

			if err := json.Unmarshal(reqBody, &dat); err != nil {
				panic(err)
			}
			enc := json.NewEncoder(os.Stdout)
			enc.Encode(dat)
		}
		defer resp.Body.Close()

		var res map[string]result

		json.NewDecoder(resp.Body).Decode(&res)
		if res != nil {
			json.NewEncoder(os.Stdout).Encode(res)
			off_projects := viper.Get("off-projects")

			if off_projects == nil {
				fmt.Println(`No off projects defined, these hours might no be correct, please check the off-project command.
				If your workspace doesn't have any off-projects, you can simple ignore this message.`)
			}
			// For each day, + or - (ignore off-hours).
			// flex hours should be subtracted
			// append hours if day already in list.
			// divide with how many days in list.
			// print saldo.
		} else {
			fmt.Println("Could not get report:", err)
		}
	},
}
