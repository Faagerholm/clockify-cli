package utils

import (
	"math/rand"
	"time"
)

func RandomExitGreeting() string {
	greetings := []string{
		"Goodbye!",
		"See you later!",
		"See you next time, buddy!",
		"See you later, alligator!",
		"See you later, buddy!",
		"Have a nice day!",
		"Byeee!",
		"Bye!",
		"May the force be with you!",
		"Adios!",
		"Auf wiedersehen!",
		"Au revoir!",
		"Tchau!",
		"Moin!",
		"Hasta luego!",
		"Jambo!",
	}
	return greetings[randomInt(0, len(greetings))]
}

func randomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())

	return min + rand.Intn(max-min)
}
