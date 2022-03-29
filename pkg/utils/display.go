package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

func DisplayMonthInCLI(date time.Time, fullMonth bool, markedDays []int, c *color.Color) {
	if c == nil {
		c = color.New(color.FgWhite)
	}

	firstDay := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.UTC)
	var lastDay time.Time
	if fullMonth {
		lastDay = time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC)
	} else {
		lastDay = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	}

	// Display month
	lineLength := 20
	padding := (lineLength - len(date.Month().String())) / 2
	fmt.Printf("\n%s %s\n", fmt.Sprint(strings.Repeat(" ", padding)), fmt.Sprint(date.Month().String()))

	// Weekday header
	fmt.Printf("%s %s %s %s %s %s %s\n",
		"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su")

	isMarked := func(day int) bool {
		for _, markedDay := range markedDays {
			if day == markedDay {
				return true
			}
		}
		return false
	}

	// Add padding to first week
	for i := 0; i < int(firstDay.Weekday()-1); i++ {
		fmt.Printf("   ")
	}

	day := firstDay
	for i := 1; i <= lastDay.Day(); i++ {
		if isMarked(day.Day()) {
			c.Printf("%2d ", day.Day())
		} else {
			fmt.Printf("%2d ", day.Day())
		}
		if day.Weekday() == time.Sunday {
			fmt.Printf("\n")
		}
		day = firstDay.AddDate(0, 0, i)
	}

	fmt.Printf("\n")
}
