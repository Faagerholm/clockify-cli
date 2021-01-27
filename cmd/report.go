package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
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

		loc, _ := time.LoadLocation("UTC")
		cur_time := time.Now().In(loc)
		end_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", cur_time.Year(), cur_time.Month(), cur_time.Day(), cur_time.Hour(), cur_time.Minute(), cur_time.Second())

		var res *Result
		res = new(Result)

		for i := 0; i >= -1; i-- {
			cur_time = cur_time.AddDate(i, 0, 0)
			first_day := fmt.Sprintf("%d-01-01T00:00:00Z", cur_time.Year())
			last_day := fmt.Sprintf("%d-12-31T00:00:00Z", cur_time.Year())
			tmp_res, _ := getEntries(first_day, last_day)
			res.addEntries(tmp_res)
		}

		entry := createEntry(res)
		if entry != nil {
			first_day := findFirstDay(entry)
			first_day_str := fmt.Sprintf("%d-%02d-%02d", first_day.Year(), first_day.Month(), first_day.Day())
			log.Println(first_day_str)
			work_days := daysCalc(end_time[0:10], first_day_str)

			//TODO: remove hard-coded part-time day
			var part_time_days float64
			part_time_percent := 1 - 0.8 // => 80%
			if true {
				part_time_days = daysCalc(end_time[0:10], "2020-11-01")
			} else {
				part_time_days = 0
			}

			log.Println(part_time_days)
			log.Println(entry.Duration)
			var balance = (float64(entry.Duration) - work_days*week_seconds + part_time_days*week_seconds*part_time_percent) / (60 * 60)
			fmt.Printf("----------------------\nYou have worked %.0f days.\nYou have worked %.2f hours, and the recommended amount is %.2f hours.\nWhich makes your balance %dh%dmin  (%.2f)\nThis calculator includes today.\n",
				work_days, // days
				float64(entry.Duration)/float64(60.0*60.0),                                              // hours worked
				float64((work_days*week_seconds-part_time_days*week_seconds*part_time_percent)/(60*60)), // recommended hours
				int(balance), // balance (hours)
				int((math.Abs(balance)-math.Abs(float64(int(balance))))*60), // balance (minutes)
				balance) // balance (decimal)
		} else {
			fmt.Println("Could not get report")
		}
	},
}

func (result *Result) addEntries(newResult *Result) []Entry {
	result.Entries = append(result.Entries, newResult.Entries...)
	return result.Entries
}

func findFirstDay(entry *Entry) time.Time {
	first_day := time.Now()
	for _, day := range entry.Children {
		d, _ := time.Parse("2006-01-02", day.Name)
		if first_day.Sub(d) > 0 {
			first_day = d
		}
	}
	return first_day
}

func createEntry(result *Result) *Entry {
	var entry *Entry
	entry = new(Entry)
	entry.Name = result.Entries[0].Name
	for _, e := range result.Entries {
		entry.Children = append(entry.Children, e.Children...)
		entry.Duration += e.Duration
	}
	log.Println(entry)
	return entry
}

func getEntries(start_time string, end_time string) (*Result, error) {
	user := viper.GetString("USER-ID")

	if len(user) == 0 {
		fmt.Println("User not set, please run `clockify user` to update the current user.")
		return nil, errors.New("User not set, please run `clockify user`")
	}

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
	res, err := requestReport(reqBody)

	return res, err
}

func requestReport(body []byte) (*Result, error) {
	key := viper.GetString("API-KEY")
	workspace := viper.GetString("WORKSPACE")

	if len(key) == 0 {
		return nil, errors.New("API KEY NOT SET! Please run clockify app-key")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("https://reports.api.clockify.me/v1/workspaces/%s/reports/summary", workspace), bytes.NewBuffer(body))
	req.Header.Set("X-API-KEY", key)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	res := new(Result)

	json.NewDecoder(resp.Body).Decode(&res)
	return res, err
}

func daysCalc(today string, start string) float64 {

	// Get first day of reports (of this year)
	first_day, _ := time.Parse("2006-01-02", start)

	// today
	now, _ := time.Parse("2006-01-02", today)
	days := (now.Sub(first_day).Hours() / 24) + 1

	// This is ugly, but works for now. Use fancy algorithm later
	var this_day = first_day
	for (this_day.Sub(now).Hours() / 24) <= 0 {
		// TODO: Compare if day == holiday
		if this_day.Weekday().String() == "Saturday" || this_day.Weekday().String() == "Sunday" {
			days -= 1
		}
		// Add one day.
		this_day = this_day.Add(time.Hour * 24)
	}
	return days
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
