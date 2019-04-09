package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Init(token string, state bool) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	bot.Debug = state

	log.Printf("Authorized on account %s,\nDebuging mode: %t", bot.Self.UserName, state)
	return bot, nil
}

func Start(b *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		userMessage := update.Message.Text
		chatIdMessage := update.Message.Chat.ID
		log.Printf("[%s] %s", update.Message.From.UserName, userMessage)

		if update.Message.IsCommand() {
			message := tgbotapi.NewMessage(chatIdMessage, "")
			switch update.Message.Command() {
			case "start":
				message.Text = "Hello!\nI calculate how many Cyrillic characters are in the articles on the medium or telegraph\nAlso..."
			default:
				message.Text = "I don't know that command"
			}
			b.Send(message)
			continue
		}
		if !IsCorrectURL(userMessage) {
			message := tgbotapi.NewMessage(chatIdMessage, "Error!\nOnly website http[s]://...")
			b.Send(message)
		} else {
			size := StartScrape(userMessage)
			message := tgbotapi.NewMessage(chatIdMessage, strconv.Itoa(size))
			b.Send(message)
		}
	}
}
