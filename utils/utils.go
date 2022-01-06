package utils

import (
	"fmt"
	"log"
	"time"
)

func GetFirstAndLastDayOfMonth() (
	firstOfMonth time.Time,
	lastOfMonth time.Time,
) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth = firstOfMonth.AddDate(0, 1, -1)
	return
}

func DisplayProjects(projectList map[string]int) {
	for p, d := range projectList {
		minutes := 0
		if (0 < d && d < 3600) || d%3600 != 0 {
			res := d % 3600
			minutes = 60 / (3600 / res)
		}
		log.Printf("%s: %dh%dmin\n", p, d/3600, minutes)
	}
}

func DisplayHoursFromMinues(minutes int) string {
	fmt.Println(minutes)
	min := 0
	if (0 < minutes && minutes < 3600) || minutes%3600 != 0 {
		res := minutes % 3600
		min = 60 / (3600 / res)
	}
	return fmt.Sprintf("%dh%dmin", int(minutes/3600), min)
}
