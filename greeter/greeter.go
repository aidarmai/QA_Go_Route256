package greeter

import (
	"fmt"
	"strings"
)

func Greet(name string, hour int) string {
	greeting := "Good night"

	if hour >= 6 && hour < 12 {
		greeting = "Good morning"
	} else if hour >= 12 && hour < 18 {
		greeting = "Hello"
	} else if hour >= 18 && hour < 22 {
		greeting = "Good evening"
	} else if hour < 0 || hour > 24 {
		greeting = "Hello"
	}
	trimmedName := strings.Trim(name, " ")
	return fmt.Sprintf("%s %s!", greeting, strings.Title(trimmedName))
}
