package main

import (
	"Character-Count-Bot/utils"
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/proxy"
)

type Config struct {
	Token string
	Proxy string
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(configuration.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		userMessage := update.Message.Text
		msgChatID := update.Message.Chat.ID
		log.Printf("[%s] %s", update.Message.From.UserName, userMessage)

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(msgChatID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = "Hello!\nI calculate how many Cyrillic characters are in the articles on the medium or telegraph"
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
			continue
		}
		if !utils.IsCorrectURL(userMessage) {
			msg := tgbotapi.NewMessage(msgChatID, "Error!\nOnly website http[s]://...")
			bot.Send(msg)
		} else {
			c := colly.NewCollector(
				colly.Async(true),
			)
			var domain string = utils.GetDomain(userMessage)
			if domain == "telegra" {
				rp, err := proxy.RoundRobinProxySwitcher(configuration.Proxy)
				if err != nil {
					text := "Error when installing proxy\n" + err.Error()
					msg := tgbotapi.NewMessage(msgChatID, text)
					bot.Send(msg)
				}
				c.SetProxyFunc(rp)
			}
			var contentPage string
			var querySelectors map[string][]string = map[string][]string{
				"medium":  {`.postArticle-content`, "section"},
				"telegra": {`.tl_article`, "article"},
			}
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
			size := strconv.Itoa(utf8.RuneCountInString(preparedText))
			msg := tgbotapi.NewMessage(msgChatID, size)
			bot.Send(msg)
		}
	}
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
