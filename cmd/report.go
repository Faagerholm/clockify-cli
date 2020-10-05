package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

/* Request type */
type report struct {
	Start         string         `json:"dateRangeStart"`
	End           string         `json:"dateRangeEnd"`
	SummaryFilter *report_filter `json:"summaryFilter,omitempty"`
	SortOrder     string         `json:"sortOrder"`
	Users         *report_user   `json:"users,omitempty"`
}
type report_filter struct {
	Groups []string `json:"groups"`
}
type report_user struct {
	Ids      []string `json:"ids"`
	Contains string   `json:"contains"`
	Status   string   `json:"status"`
}

/* Response types */
type Result struct {
	Entries []Entry `json:"groupOne"`
}

type Entry struct {
	Name     string
	Duration int
	Children []groupChild
}
type groupChild struct {
	Duration int
	Name     string
}

/* const */
const (
	week_seconds = 27000 // 37.5 * 60 * 60
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Display if you're above or below zero balance.",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.GetString("API-KEY")
		workspace := viper.GetString("WORKSPACE")
		user := viper.GetString("USER-ID")

		if len(key) == 0 {
			fmt.Println("API KEY NOT SET! Please run clockify app-key")
			return
		}

		if len(user) == 0 || len(workspace) == 0 {
			fmt.Println("User not set, please run `clockify user` to update the current user.")
			return
		}

		loc, _ := time.LoadLocation("UTC")
		cur_time := time.Now().In(loc)

		start_time := fmt.Sprintf("%d-01-01T00:00:00Z", cur_time.Year())
		end_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", cur_time.Year(), cur_time.Month(), cur_time.Day(), cur_time.Hour(), cur_time.Minute(), cur_time.Second())

		users := report_user{
			Ids:      []string{user},
			Contains: "CONTAINS",
			Status:   "All",
		}
		filter := report_filter{
			Groups: []string{"user", "date"},
		}
		reqBody, err := json.Marshal(report{
			Start:         start_time,
			End:           end_time,
			SummaryFilter: &filter,
			SortOrder:     "Ascending",
			Users:         &users,
		})
		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		req, _ := http.NewRequest("POST", fmt.Sprintf("https://reports.api.clockify.me/v1/workspaces/%s/reports/summary", workspace), bytes.NewBuffer(reqBody))
		req.Header.Set("X-API-KEY", key)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		res := new(Result)

		json.NewDecoder(resp.Body).Decode(&res)
		if res != nil {
			// json.NewEncoder(os.Stdout).Encode(res)

			// Get first day of reports (of this year)
			first_day, _ := time.Parse("2006-01-02", res.Entries[0].Children[0].Name)
			// today
			today, _ := time.Parse("2006-01-02", end_time[0:10])
			max_days := (today.Sub(first_day).Hours() / 24) + 1

			// This is ugly, but works for now. Use fancy algorithm later
			var this_day = first_day
			for (this_day.Sub(today).Hours() / 24) <= 0 {
				// TODO: Compare if day == holiday
				if this_day.Weekday().String() == "Saturday" || this_day.Weekday().String() == "Sunday" {
					max_days -= 1
				}
				// Add one day.
				this_day = this_day.Add(time.Hour * 24)
			}
			var balance = (float64(res.Entries[0].Duration) - max_days*week_seconds) / (60 * 60)
			fmt.Printf("----------------------\nYou have worked %.0f days.\n", max_days)
			fmt.Printf("You have worked %.2f hours, and the recommended amount is %.2f hours.\nWhich makes your balance %dh%dmin  (%.2f)\n",
				float64(res.Entries[0].Duration)/float64(60.0*60.0),
				float64(max_days*week_seconds/(60*60)),
				int(balance),
				int((math.Abs(balance)-math.Abs(float64(int(balance))))*60),
				balance)

		} else {
			fmt.Println("Could not get report:", err)
		}
	},
}

func stringInSlice(a string, list interface{}) bool {

	listSlice, ok := list.([]interface{})
	if !ok {
		return false
	}
	for _, v := range listSlice {
		if a == v {
			return true
		}
	}
	return false
}
