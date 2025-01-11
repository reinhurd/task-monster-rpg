package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"rpgMonster/internal/core"
	"rpgMonster/internal/model"
)

var lastChatID int64

type TGBot struct {
	bot *tgbotapi.BotAPI
	u   tgbotapi.UpdateConfig
	svc *core.Service
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

			resp = "Hello, " + update.Message.From.UserName + "!" + " You said: " + update.Message.Text + ", to get help type" + model.HELP
			userTelegramID := update.Message.From.ID

			switch {
			case strings.Contains(update.Message.Text, model.CREATE_TASK_GPT):
				userID, err := t.svc.ValidateUserTG(int(userTelegramID)) //todo think about int and int64 in tgID
				if err != nil {
					resp = err.Error()
					break
				}
				splStr := strings.Split(update.Message.Text, " ")
				if len(splStr) < 2 {
					resp = "Please specify request"
				} else {
					task, err := t.svc.CreateTaskFromGPTByRequest(splStr[1], userID)
					if err != nil {
						resp = err.Error()
					}
					resp = fmt.Sprintf(model.Commands[model.CREATE_TASK_GPT], task)
				}
			case strings.Contains(update.Message.Text, model.CONNECT_USER):
				userID, err := t.svc.ValidateUserTG(int(userTelegramID))
				if err != nil {
					resp = err.Error()
					break
				}
				spStr := strings.Split(update.Message.Text, " ")
				if len(spStr) < 2 {
					resp = "Please specify user ID"
				} else {
					userID, err = t.svc.ConnectUserToTG(spStr[1], int(userTelegramID))
					if err != nil {
						resp = err.Error()
					}
					resp = fmt.Sprintf(model.Commands[model.CONNECT_USER], userID)
				}
			case strings.Contains(update.Message.Text, model.HELP):
				resp = model.Commands[model.HELP]
			}

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

func StartBot(tgToken string, isDebug bool, svc *core.Service) (*TGBot, error) {
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
		svc: svc,
	}

	return tgbot, nil
}
