package tgbot

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func loadTokenFromEnv() (string, error) {
	token, ok := os.LookupEnv("TG_API_SECRET_KEY")
	if !ok {
		return "", fmt.Errorf("environment variable OPENAI_API_SECRET_KEY not set")
	}
	return token, nil
}

func StartBot() {
	privateToken, err := loadTokenFromEnv()
	if err != nil {
		log.Fatalf("error in loadTokenFromEnv: %v", err)
	}
	bot, err := tgbotapi.NewBotAPI(privateToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var resp string
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			resp = "Hello, " + update.Message.From.UserName + "!" + " You said: " + update.Message.Text

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
