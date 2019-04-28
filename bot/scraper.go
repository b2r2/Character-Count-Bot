package bot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/gocolly/colly"
)

func GetCountSymbolsInArticle(userMessage string, c *Config) (int, error) {
	var content string
	var err error
	if comp := regexp.MustCompile(c.Site.Domain); comp.MatchString(userMessage) {
		content, err = scrapeSite(userMessage, c)
		if err != nil {
			return 0, err
		}
	} else {
		content, err = scrapeMedium(c.Medium, userMessage)
		if err != nil {
			return 0, err
		}
	}
	text, err := parse(content)
	if err != nil {
		return 0, err
	}
	size := utf8.RuneCountInString(text)
	return size, nil
}

func scrapeMedium(domain, url string) (string, error) {
	col := colly.NewCollector(
		colly.Async(true),
	)
	if comp := regexp.MustCompile("edit$"); comp.MatchString(url) {
		url = comp.ReplaceAllString(url, "")
	}
	var querySelectors map[string][]string = map[string][]string{
		domain: {`.postArticle-content`, "section"},
	}
	var text string
	var querySelector string = querySelectors[domain][0]
	col.OnHTML(querySelector, func(e *colly.HTMLElement) {
		var tag string = querySelectors[domain][1]
		text = e.ChildText(tag)
	})
	col.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})
	col.Visit(url)
	col.Wait()
	return text, nil
}

func scrapeSite(url string, c *Config) (string, error) {
	var wpr WPResponse
	var postNumber string
	re := regexp.MustCompile(`[0-9]+`)
	if re.MatchString(url) {
		postNumber = string(re.Find([]byte(url)))
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", c.Site.URL+postNumber, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(c.Site.Login, c.Site.Password)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	err = json.Unmarshal(data, &wpr)
	if err != nil {
		return "", err
	}
	return wpr.Content.Rendered, nil
}

func parse(text string) (string, error) {
	re, err := regexp.Compile("\\p{Cyrillic}")
	if err != nil {
		return "", err
	}
	temp := re.FindAllString(text, -1)
	var total string
	for _, t := range temp {
		total += t
	}
	return total, nil
}
