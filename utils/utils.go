package utils

import (
	"strings"
)

func IsCorrectURL(url string) bool {
	correctUrl := []string{
		"medium",
	}
	state := false
	line := strings.Split(url, "://")
	if len(line) == 2 {
		line = strings.Split(line[1], ".")
		if line[0] == correctUrl[0] {
			state = true
		}
	}
	return state
}

func GetDomain(url string) string {
	line := strings.Split(url, "://")
	line = strings.Split(line[1], ".")
	return line[0]
}
