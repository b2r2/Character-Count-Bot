package bot

import (
	"fmt"
	"log"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
)

func StartScrape(userMessage string, c *Config) int {
	col := colly.NewCollector(
		colly.Async(true),
	)
	var comp = regexp.MustCompile("edit$")
	if comp.MatchString(userMessage) {
		userMessage = comp.ReplaceAllString(userMessage, "")
	}
	comp = regexp.MustCompile(c.Scraping.Site)
	if comp.MatchString(userMessage) {
		userMessage = fmt.Sprintf(userMessage + "?no_cache")
	}
	var querySelectors map[string][]string = map[string][]string{
		c.Scraping.Medium: {`.postArticle-content`, "section"},
		c.Scraping.Site: {fmt.Sprintf(`.post-%s`, func(s string) string {
			re := regexp.MustCompile(`[0-9]+`)
			return re.FindAllString(s, -1)[0]
		}(userMessage)), `.td-post-content`},
	}
	var domain string = GetDomain(userMessage)
	var querySelector string = querySelectors[domain][0]
	var contentPage string
	col.OnHTML(querySelector, func(e *colly.HTMLElement) {
		var tag string = querySelectors[domain][1]
		contentPage = e.ChildText(tag)
	})
	col.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	col.Visit(userMessage)
	col.Wait()
	preparedText, err := parsePage(contentPage)
	if err != nil {
		log.Fatal(err)
	}
	size := utf8.RuneCountInString(preparedText)
	return size
}

func parsePage(contentPage string) (string, error) {
	re, err := regexp.Compile("\\p{Cyrillic}")
	if err != nil {
		return "", err
	}
	temp := re.FindAllString(contentPage, -1)
	var totString string
	for _, t := range temp {
		totString += t
	}
	return totString, nil
}
