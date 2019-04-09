package bot

import (
	"log"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
)

func StartScrape(userMessage string) int {
	c := colly.NewCollector(
		colly.Async(true),
	)
	var comp = regexp.MustCompile("edit$")
	if comp.MatchString(userMessage) {
		userMessage = comp.ReplaceAllString(userMessage, "")
	}
	var contentPage string
	var querySelectors map[string][]string = map[string][]string{
		"medium":  {`.postArticle-content`, "section"},
		"telegra": {`.tl_article`, "article"},
	}
	var domain string = GetDomain(userMessage)
	var querySelector string = querySelectors[domain][0]
	c.OnHTML(querySelector, func(e *colly.HTMLElement) {
		var tag string = querySelectors[domain][1]
		contentPage = e.ChildText(tag)
	})
	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	c.Visit(userMessage)
	c.Wait()
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
