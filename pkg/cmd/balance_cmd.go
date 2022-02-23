package cmd

import (
	"fmt"
	"log"
	"time"

	api "github.com/Faagerholm/clockify-cli/pkg/API"
	model "github.com/Faagerholm/clockify-cli/pkg/Model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var CheckBalanceCmd = &cobra.Command{
	Use:   "check-balance",
	Short: "Check balance",
	Long:  `Check the balance of the current account`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckBalance()
	},
}

var AddPartTimeCmd = &cobra.Command{
	Use:   "add-part-time",
	Short: "Add part-time work to your account",
	Long:  `Add part-time, if you have been working less during a time period.`,
	Run: func(cmd *cobra.Command, args []string) {
		AddPartTimeTimespan()
	},
}

func CheckBalance() {
	user, err := api.GetAllEntries()

	if err != nil {
		log.Println(err)
		return
	}
	if err != nil {
		log.Println(err)
		return
	}

	first_day := extractFirstDayFromEntries(user.Entries)
	expected_work, _ := countExpectedWorkingTime(first_day, time.Now())

	balance := countBalance(user.Duration, expected_work)

	fmt.Printf("Balance: %.2f hours\n", balance)
}

func AddPartTimeTimespan() {
	var partTimes []model.PartTime
	viper.UnmarshalKey("part-time", &partTimes)

	var newPartTime model.PartTime
	fmt.Print("Enter start date (YYYY-MM-DD): ")
	fmt.Scanln(&newPartTime.Start)
	fmt.Print("Enter end date (YYYY-MM-DD): ")
	fmt.Scanln(&newPartTime.End)
	fmt.Print("Enter capacity (0-100): ")
	fmt.Scanln(&newPartTime.Capacity)

	partTimes = append(partTimes, newPartTime)
	viper.Set("part-time", partTimes)
	viper.WriteConfig()

	fmt.Printf("Added part-time: %v\n", newPartTime)
}

func extractFirstDayFromEntries(entries []model.ReportEntry) time.Time {
	first_day := time.Now()
	for _, day := range entries {
		d, _ := time.Parse("2006-01-02", day.Date)
		if first_day.Sub(d) > 0 {
			first_day = d
		}
	}
	return first_day
}

func countWorkingDaysBetweenTwoDays(start, end time.Time) int32 {

	days := int32(end.Sub(start).Hours()/24) + 1

	// This is ugly, but works for now. Use fancy algorithm later
	// Maybe some AI >_<
	var this_day = start
	for (this_day.Sub(end).Hours() / 24) <= 0 {
		// We are marking holidays in clockify, so we want to include them
		if this_day.Weekday().String() == "Saturday" || this_day.Weekday().String() == "Sunday" {
			days -= 1
		}
		// Add one day.
		this_day = this_day.Add(time.Hour * 24)
	}
	return days
}

func countExpectedWorkingTime(start, end time.Time) (int64, int32) {
	working_days := countWorkingDaysBetweenTwoDays(start, end)
	expected_working_days_in_seconds := int64(float64(working_days) * 7.5 * 60.0 * 60.0)
	expected_working_days_in_seconds = subtractPartTimeWork(expected_working_days_in_seconds)
	return expected_working_days_in_seconds, working_days
}

func subtractPartTimeWork(duration int64) int64 {
	var partTimes []model.PartTime
	viper.UnmarshalKey("part-time", &partTimes)

	for _, partTime := range partTimes {
		start, err := time.Parse("2006-01-02", partTime.Start)
		if err != nil {
			log.Printf("Error parsing start time of part-time: %v\n", err)
			return duration
		}
		end, err := time.Parse("2006-01-02", partTime.End)
		if err != nil {
			log.Printf("Error parsing end date of part-time: %v\n", err)
			return duration
		}
		days := countWorkingDaysBetweenTwoDays(start, end)
		reduction := int64(float64(days) * (7.5 * float64(100-partTime.Capacity) / 100) * 60.0 * 60.0)
		duration -= reduction
	}
	return duration
}

// Count balance, duration is in seconds
func countBalance(duration int64, expected_working_days_in_seconds int64) float64 {
	balance_in_seconds := duration - expected_working_days_in_seconds
	balance_in_hours := float64(balance_in_seconds) / 3600.0
	return balance_in_hours
}
