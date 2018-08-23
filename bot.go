package main

import (
	"Character-Count-Bot/config"
	"Character-Count-Bot/utils"
	"log"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gocolly/colly"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.Configuration["TOKEN"])
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		userName := update.Message.From.UserName
		userMessage := update.Message.Text
		msgChatID := update.Message.Chat.ID
		log.Printf("[%s] %s", userName, userMessage)

		if !utils.IsCorrectURL(userMessage) {
			msg := tgbotapi.NewMessage(msgChatID, "Error!\nOnly website http[s]://medium.com/...")
			bot.Send(msg)
		} else {

			c := colly.NewCollector(
				colly.Async(true),
			)

			var contentPage string
			var querySelector = `.postArticle-content`

			c.OnHTML(querySelector, func(e *colly.HTMLElement) {
				var tag string = "section"
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
