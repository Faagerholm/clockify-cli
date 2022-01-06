package clockify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Faagerholm/clockify-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var BalanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Display if you're above or below zero balance.",
	Run: func(cmd *cobra.Command, args []string) {
		balance, err := getTotalBalance()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Balance: ", utils.DisplayHoursFromMinues(int(balance)))
	},
}

var WorkLogCmd = &cobra.Command{
	Use:   "c-report",
	Short: "Get this months report",
	Run: func(cmd *cobra.Command, args []string) {
		firstDayOfMonth, lastDayOfMonth := utils.GetFirstAndLastDayOfMonth()
		worklog, err := GetWorklog(firstDayOfMonth, lastDayOfMonth)
		if err != nil {
			log.Fatal(err)
		}

		projects := make(map[string]int)

		for _, e := range worklog.Entries {
			for _, c := range e.Children {
				for _, p := range c.Children {

					_, ok := projects[p.Project]
					if ok {
						projects[p.Project] += p.Duration
					} else {
						projects[p.Project] = p.Duration
					}
				}
			}
		}
		displayProjects(projects)
	},
}

func GetWorklog(start_time time.Time, end_time time.Time) (*ClockifyReport, error) {
	user := viper.GetString("USER-ID")

	if len(user) == 0 {
		log.Println("User not set, please run `clockify user` to update the current user.")
		return nil, errors.New("user not set, please run `clockify user`")
	}

	reqBody, err := json.Marshal(Report{
		Start:         start_time.Format("2006-01-02T15:04:05Z"),
		End:           end_time.Format("2006-01-02T15:04:05Z"),
		SummaryFilter: &report_filter{[]string{"user", "date", "project"}},
		SortOrder:     "Ascending",
		Users:         &report_user{Ids: []string{user}, Contains: "CONTAINS", Status: "ALL"},
	})

	if err != nil {
		log.Fatal(err)
	}
	res, err := requestReport(reqBody)

	return res, err
}

func requestReport(body []byte) (*ClockifyReport, error) {
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

	res := new(ClockifyReport)
	json.NewDecoder(resp.Body).Decode(&res)
	return res, err
}

func daysCalc(end string, start string) int {

	// Parse the date inputs
	first_day, _ := time.Parse("2006-01-02", start)
	last_day, _ := time.Parse("2006-01-02", end)

	// Get the number of days between the two dates
	days := int((last_day.Sub(first_day).Hours() / 24) + 1)

	// This is ugly, but works for now. Use fancy algorithm later
	var this_day = first_day
	for (this_day.Sub(last_day).Hours() / 24) <= 0 {
		if this_day.Weekday().String() == "Saturday" || this_day.Weekday().String() == "Sunday" {
			days -= 1
		}
		// Increment the day
		this_day = this_day.Add(time.Hour * 24)
	}

	log.Printf("Accumulated days: %d\nBetween: %s - %s", days, start, end)
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

func displayProjects(projectList map[string]int) {
	for p, d := range projectList {
		minutes := 0
		if (0 < d && d < 3600) || d%3600 != 0 {
			res := d % 3600
			minutes = 60 / (3600 / res)
		}
		log.Printf("%s: %dh%dmin\n", p, d/3600, minutes)
	}
}

func getTotalBalance() (float64, error) {
	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()

	total_hours := 0.0
	days := 0

	// Partial year
	worklog, err := GetWorklog(time.Date(
		currentYear, time.January, 1, 0, 0, 0, 0, now.Location()),
		time.Date(currentYear, currentMonth, currentDay+1, 0, 0, 0, 0, now.Location()))

	if err != nil {
		return 0, err
	}
	total_hours += float64(worklog.Entries[0].Total)
	days += daysCalc(
		fmt.Sprintf("%d-%d-%d", currentYear, currentMonth, currentDay),
		worklog.Entries[0].Children[0].Date)
	log.Printf("%d: %d\n", currentYear, worklog.Entries[0].Total)
	log.Printf("You worked %f hours in year %d\n", float32(worklog.Entries[0].Total/60/60), currentYear)

	// Previous years -- until no further years are found
	for {

		currentYear--
		worklog, err := GetWorklog(time.Date(currentYear, time.January, 1, 0, 0, 0, 0, now.Location()), time.Date(currentYear, time.December, 31, 0, 0, 0, 0, now.Location()))
		if err != nil {
			return 0, err
		}
		if len(worklog.Entries) == 0 {
			break
		}
		log.Printf("%d: %d\n", currentYear, worklog.Entries[0].Total)
		days += daysCalc(
			fmt.Sprintf("%d-%d-%d", currentYear, time.December, 31),
			worklog.Entries[0].Children[0].Date)
		log.Printf("You worked %f hours in year %d\n", float32(worklog.Entries[0].Total)/60.0/60.0, currentYear)
		total_hours += float64(worklog.Entries[0].Total)
	}
	log.Println("You worked", float32(total_hours)/60.0/60.0, "hours in", days, "days")
	balance := float64(total_hours/60.0) - (float64(days) * 7.5)
	balance += accountForPartialDays()
	return balance, nil
}

type PartTime struct {
	StartDate string
	EndDate   string
	Capacity  float64
}

func accountForPartialDays() float64 {
	// Check how to account for multiple partial day intervals

	partialTime := &PartTime{
		StartDate: viper.GetString("part-time.start-date"),
		EndDate:   viper.GetString("part-time.end-date"),
		Capacity:  viper.GetFloat64("part-time.capacity"),
	}
	if partialTime == nil {
		return 0
	}

	days := daysCalc(partialTime.EndDate, partialTime.StartDate)
	log.Printf("Partial days: %d\n", days)
	return (float64(days) * ((7.5) - (7.5 * partialTime.Capacity))) * 60.0
}
