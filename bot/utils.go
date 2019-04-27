package bot

import (
	"strings"
)

func IsCorrectURL(url string, c *Config) bool {
	correctUrl := []string{
		c.Scraping.Medium,
		c.Scraping.Site,
	}
	state := false
	line := strings.Split(url, "://")
	if len(line) == 2 {
		line = strings.Split(line[1], ".")
		for _, curl := range correctUrl {
			if line[0] == curl {
				state = true
				break
			}
		}
	}
	return state
}

func GetDomain(url string) string {
	line := strings.Split(url, "://")
	line = strings.Split(line[1], ".")
	return line[0]
}
