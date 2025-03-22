package telegram

import (
	"context"
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
	for update := range updates {
		var resp string
		if update.Message != nil { // If we got a message
			log.Info().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			userTelegramID := update.Message.From.ID

			//check if user exists and show his current tasks
			userID, err := t.svc.ValidateUserTG(userTelegramID)
			if userID == "" {
				//message that user need to register
				resp = "Please register first, type " + model.Commands[model.CREATE_USER] + " <login> <password>"
			} else if err != nil {
				resp = err.Error()
			} else {
				tasks, err := t.svc.GetListTasksByUserID(context.Background(), userID)
				if err != nil {
					resp = err.Error()
				} else {
					for _, task := range tasks {
						resp = fmt.Sprintf(model.Commands[model.TASK_LIST], task)
					}
				}
			}

			switch {
			case strings.Contains(update.Message.Text, model.CREATE_TASK_GPT):
				userID, err := t.svc.ValidateUserTG(userTelegramID)
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
					} else {
						resp = fmt.Sprintf(model.Commands[model.CREATE_TASK_GPT], task)
					}
				}
			case strings.Contains(update.Message.Text, model.CONNECT_USER):
				spStr := strings.Split(update.Message.Text, " ")
				if len(spStr) < 3 {
					resp = "Please specify user login and password"
				} else {
					userID, _, err := t.svc.CheckPassword(spStr[1], spStr[2])
					if err != nil {
						return err
					}
					err = t.svc.ConnectUserToTG(userID, userTelegramID)
					if err != nil {
						resp = err.Error()
					} else {
						resp = fmt.Sprintf(model.Commands[model.CONNECT_USER], userID)
					}
				}
			case strings.Contains(update.Message.Text, model.TASK_LIST):
				userID, err := t.svc.ValidateUserTG(userTelegramID)
				if err != nil {
					resp = err.Error()
					break
				}
				tasks, err := t.svc.GetListTasksByUserID(context.Background(), userID)
				if err != nil {
					resp = err.Error()
				} else {
					for _, task := range tasks {
						resp += fmt.Sprintf(model.Commands[model.TASK_LIST], task)
					}
				}
			case strings.Contains(update.Message.Text, model.CREATE_TASK_GPT):
				userID, err := t.svc.ValidateUserTG(userTelegramID)
				if err != nil {
					resp = err.Error()
					break
				}
				splStr := strings.Split(update.Message.Text, " ")
				if len(splStr) < 3 {
					resp = "Please specify request"
				} else {
					var task model.Task
					task.Title = splStr[1]
					task.Description = splStr[2]
					task.Executor = userID
					err = t.svc.CreateTask(context.Background(), &task)
					if err != nil {
						resp = err.Error()
					} else {
						resp = fmt.Sprintf(model.Commands[model.CREATE_TASK], task)
					}
				}
			case strings.Contains(update.Message.Text, model.UPDATE_TASK):
				userID, err := t.svc.ValidateUserTG(userTelegramID)
				if err != nil {
					resp = err.Error()
					break
				}
				splStr := strings.Split(update.Message.Text, " ")
				if len(splStr) < 4 {
					resp = "Please specify task ID, goal and description"
				} else {
					var task model.Task
					task.BizId = splStr[1]
					task.Title = splStr[2]
					task.Description = splStr[3]
					task.Executor = userID
					err = t.svc.UpdateTask(context.Background(), &task)
					if err != nil {
						resp = err.Error()
					} else {
						resp = fmt.Sprintf(model.Commands[model.UPDATE_TASK], task)
					}
				}
			case strings.Contains(update.Message.Text, model.CREATE_USER):
				spStr := strings.Split(update.Message.Text, " ")
				if len(spStr) < 3 {
					resp = "Please specify login and password"
				} else {
					userID, err := t.svc.CreateUserFromTG(spStr[1], spStr[2], userTelegramID)
					if err != nil {
						resp = err.Error()
					} else {
						resp = fmt.Sprintf(model.Commands[model.CREATE_USER], userID)
					}
				}
			case strings.Contains(update.Message.Text, model.START):
				if userID == "" {
					resp = "Please register first, type " + model.Commands[model.CREATE_USER] + " <login> <password> or connect user " + model.Commands[model.CONNECT_USER] + " <login> <password>"
				} else {
					resp = "Hello, " + update.Message.From.UserName + "!" + " You already registered as " + userID + ", to get help type " + model.HELP
				}
			case strings.Contains(update.Message.Text, model.HELP):
				resp = model.Commands[model.HELP]
			default:
				resp = "Hello, " + update.Message.From.UserName + "!" + " You said: " + update.Message.Text + ", to get help type " + model.HELP
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
