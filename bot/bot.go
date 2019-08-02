package bot

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Start(state bool, configFile string) {
	var c Config
	c.LoadScrapingConfiguration(configFile)

	b, err := tgbotapi.NewBotAPI(c.Token)
	if err != nil {
		log.Panic(err)
	}
	b.Debug = state

	log.Printf("Authorized on account %s,\nDebuging mode: %t", b.Self.UserName, state)
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
		user, err := b.GetChatMember(tgbotapi.ChatConfigWithUser{
			ChatID: c.ChatID,
			UserID: update.Message.From.ID})
		if err != nil {
			log.Panic(err)
		}
		if !(user.IsCreator() || user.IsAdministrator() || user.IsMember()) {
			continue
		}

		userMessage := update.Message.Text
		userChatID := int64(update.Message.From.ID)
		if update.Message.IsCommand() {
			message := tgbotapi.NewMessage(userChatID, "")
			switch update.Message.Command() {
			case "start":
				message.Text = "Hello!\nI calculate how many Cyrillic characters are in the articles on the medium or telegraph\nAlso..."
			default:
				message.Text = "I don't know that command"
			}
			b.Send(message)
			continue
		}
		if !IsCorrectURL(userMessage, &c) {
			message := tgbotapi.NewMessage(userChatID, "Bad request: use URL with procotol https or http")
			b.Send(message)
			continue
		}
		message := tgbotapi.NewMessage(userChatID, "")
		if size, err := GetCountSymbolsInArticle(userMessage, &c); err != nil {
			message.Text = fmt.Sprintf("Something was wrong:\n%v", err)
		} else {
			message.Text = strconv.Itoa(size)
		}
		b.Send(message)
	}
}
