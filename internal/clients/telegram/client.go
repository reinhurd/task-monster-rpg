package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

var lastChatID int64

type TGBot struct {
	bot *tgbotapi.BotAPI
	u   tgbotapi.UpdateConfig
}

func (t *TGBot) GetUpdatesChan() tgbotapi.UpdatesChannel {
	return t.bot.GetUpdatesChan(t.u)
}

func (t *TGBot) Send(chatID int64, messageID int, resp string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, resp)
	msg.ReplyToMessageID = messageID
	return t.bot.Send(msg)
}

func (t *TGBot) SendToLastChat(resp string) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(lastChatID, resp)
	msg.ReplyToMessageID = 0
	return t.bot.Send(msg)
}

func (t *TGBot) HandleUpdate(updates tgbotapi.UpdatesChannel) error {
	var err error

	for update := range updates {
		var resp string
		if update.Message != nil { // If we got a message
			log.Info().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			resp = "Hello, " + update.Message.From.UserName + "!" + " You said: " + update.Message.Text

			//todo do some logic in switch for parsing commands

			lastChatID = update.Message.Chat.ID

			if len(resp) > 4095 {
				// split message and send in parts
				for len(resp) > 4095 {
					_, err = t.Send(update.Message.Chat.ID, update.Message.MessageID, resp[:4095])
					if err != nil {
						log.Err(err).Msg("send error")
					}
					resp = resp[4095:]
				}
			}
			_, err = t.Send(update.Message.Chat.ID, update.Message.MessageID, resp)
			if err != nil {
				log.Err(err).Msg("send error")
			}
		}
	}
	return nil
}

func StartBot(tgToken string, isDebug bool) (*TGBot, error) {
	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Fatal().Err(err).Msg("tgbotapi.NewBotAPI doesn't start")
	}

	bot.Debug = isDebug

	log.Info().Msgf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	tgbot := &TGBot{
		bot: bot,
		u:   u,
	}

	return tgbot, nil
}
