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

type Result struct {
	Entries []Entry             `json:"timeentries"`
	Totals  []map[string]string `json:"totals"`
}

type Entry struct {
	ProjectId    string
	TimeInterval Interval
}
type Interval struct {
	Duration int
	Start    string
	End      string
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

const (
	week_seconds = 27000 // 37.5 * 60 * 60
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Display if you're above or below zero balance.",
	Run: func(cmd *cobra.Command, args []string) {
		key := viper.Get("API-KEY").(string)
		workspace := viper.Get("WORKSPACE")
		user := viper.GetString("USER-ID")

		loc, _ := time.LoadLocation("UTC")
		cur_time := time.Now().In(loc)

		start_time := fmt.Sprintf("%d-01-01T00:00:00Z", cur_time.Year())
		end_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", cur_time.Year(), cur_time.Month(), cur_time.Day(), cur_time.Hour(), cur_time.Minute(), cur_time.Second())

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
		}

		defer resp.Body.Close()

		res := new(Result)

		json.NewDecoder(resp.Body).Decode(&res)
		if res != nil {
			// json.NewEncoder(os.Stdout).Encode(res)
			off_projects := viper.Get("off-projects")

			if off_projects == nil {
				fmt.Println(`No off projects defined, these hours might no be correct, please check the off-project command.
				If your workspace doesn't have any off-projects, you can simple ignore this message.`)
			}
			var days []string
			var max_days = 0.0
			var duration int
			for _, v := range res.Entries {
				duration += v.TimeInterval.Duration
				if !stringInSlice(v.ProjectId, off_projects) {
					if !stringInSlice(v.TimeInterval.Start[0:10], days) {
						days = append(days, v.TimeInterval.Start[0:10])
					}
				}
			}

			if len(days) >= 2 {
				first_day, _ := time.Parse("2006-01-02", days[len(days)-1])
				last_day, _ := time.Parse("2006-01-02", days[0])
				max_days = (last_day.Sub(first_day).Hours() / 24) + 1

				// This is ugly, but works for now. Use fancy algorithm later
				var this_day = first_day
				for (this_day.Sub(last_day).Hours() / 24) <= 0 {
					// TODO: Compare if day == holiday
					if this_day.Weekday().String() == "Saturday" || this_day.Weekday().String() == "Sunday" {
						max_days -= 1
					}
					// Add one day.
					this_day = this_day.Add(time.Hour * 24)
				}
			}
			var balance = (float64(duration) - max_days*week_seconds) / (60 * 60)
			fmt.Printf("----------------------\nYou have worked %.0f days.\n", max_days)
			fmt.Printf("You have worked %.2f hours, and the recommended amount is %.2f hours.\nWhich makes your balance %dh%dmin  (%.2f)\n",
				float64(duration)/float64(60.0*60.0),
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
