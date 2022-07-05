package cmd

import (
	"fmt"
	"time"

	api "github.com/Faagerholm/clockify-cli/API"
	utils "github.com/Faagerholm/clockify-cli/Utils"
	domain "github.com/Faagerholm/clockify-cli/domain"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
)

var VerfiyMonthCmd = &cobra.Command{
	Use:   "verify-month",
	Short: "Verify month",
	Long:  `Verify the current month for missing entries`,
	Run: func(cmd *cobra.Command, args []string) {
		VerifyFullMonth()
	},
}

func VerifyFullMonth() {
	now := time.Now()
	firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC)

	verifyMonth(now, firstDay, lastDay, true)
}

func verifyMonth(now, firstDay, lastDay time.Time, fullMonth bool) {

	user, err := api.GetUserEntries(firstDay, lastDay)
	if err != nil {
		panic(err)
	}

	// loop over all days in the month
	maxHours := 0.0
	missingDays := []int{}
	for day := firstDay; day.Before(lastDay.AddDate(0, 0, 1)); day = day.AddDate(0, 0, 1) {
		// check if current day is weekday
		if day.Weekday() != time.Sunday && day.Weekday() != time.Saturday {
			// check if current day is in user entries
			if !isDayInEntries(day, user.Entries) {
				missingDays = append(missingDays, day.Day())
			}
			maxHours += 7.5
		}
	}

	if len(missingDays) != 0 {
		utils.DisplayMonthInCLI(now, fullMonth, missingDays, color.New(color.BgHiRed, color.Bold))
	} else {
		color.Green("All days in the month are present, good to go!", color.Bold)

		utils.DisplayMonthInCLI(firstDay, fullMonth, []int{}, nil)
	}

	// Count total hours for the month
	totalHours := 0.0
	for _, entry := range user.Entries {
		totalHours += float64(entry.Duration) / 3600
	}

	fmt.Printf("\nTotal hours this month: %.2f (a normal month is %.2fh) \n", totalHours, maxHours)
	fmt.Printf("This months saldo is: %.2fh\n", totalHours-maxHours)

}

func isDayInEntries(date time.Time, entries []domain.ReportEntry) bool {
	dateStr := date.Format("2006-01-02")

	for _, entry := range entries {
		if entry.Date == dateStr {
			return true
		}
	}
	return false
}
