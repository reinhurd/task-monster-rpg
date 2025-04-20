package telegram

import (
	"context"
	"fmt"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
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
			}

			lastChatID = update.Message.Chat.ID
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

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
					request := strings.Join(splStr[1:], " ")
					task, err := t.svc.CreateTaskFromGPTByRequest(request, userID)
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
					if err == mongo.ErrNoDocuments {
						resp = "User with this login and password doesn't exist"
						break
					}
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
					//sort all tasks by deadline, and put them to today, week and other slices
					todayTasks := make([]model.Task, 0)
					weekTasks := make([]model.Task, 0)
					otherTasks := make([]model.Task, 0)
					sort.Slice(tasks, func(i, j int) bool {
						return tasks[i].Deadline.Before(tasks[j].Deadline)
					})
					for _, task := range tasks {
						if task.Deadline.IsZero() {
							otherTasks = append(otherTasks, task)
						} else {
							if task.Deadline.Day() == 0 {
								todayTasks = append(todayTasks, task)
							} else if task.Deadline.Day() < 7 {
								weekTasks = append(weekTasks, task)
							}
						}
					}
					buttons := make([]tgbotapi.KeyboardButton, 0)
					resp += "---TODAY---\n"
					for _, task := range todayTasks {
						if task.BizId != "" {
							resp += fmt.Sprintf("Task description: %v\n Task %v \n", task.Description, task.BizId)
							button := tgbotapi.NewKeyboardButton(model.VIEW_TASK + " " + task.BizId)
							buttons = append(buttons, button)
						}
					}
					resp += "---WEEK---\n"
					for _, task := range weekTasks {
						if task.BizId != "" {
							resp += fmt.Sprintf("Task description: %v\n Task %v \n", task.Description, task.BizId)
							button := tgbotapi.NewKeyboardButton(model.VIEW_TASK + " " + task.BizId)
							buttons = append(buttons, button)
						}
					}
					resp += "---OTHER---\n"
					for _, task := range otherTasks {
						if task.BizId != "" {
							resp += fmt.Sprintf("Task description: %v\n Task %v \n", task.Description, task.BizId)
							button := tgbotapi.NewKeyboardButton(model.VIEW_TASK + " " + task.BizId)
							buttons = append(buttons, button)
						}
					}

					keyboard := tgbotapi.NewOneTimeReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							buttons...,
						),
					)
					msg.ReplyMarkup = keyboard
				}
			case strings.Contains(update.Message.Text, model.VIEW_TASK):
				userID, err := t.svc.ValidateUserTG(userTelegramID)
				if err != nil {
					resp = err.Error()
					break
				}
				taskID := strings.Split(update.Message.Text, " ")
				if len(taskID) < 2 {
					resp = "Please specify task ID"
				} else {
					buttons := make([]tgbotapi.KeyboardButton, 0)
					task, err := t.svc.GetTask(context.Background(), taskID[1], userID)
					if err != nil {
						resp = err.Error()
					} else {
						resp = fmt.Sprintf(model.Commands[model.VIEW_TASK], task.BizId, task.Title, task.Description, task.Completed, task.Executor, task.Reviewer, task.Deadline, task.CreatedAt, task.UpdatedAt)
						button := tgbotapi.NewKeyboardButton(model.UPDATE_TASK + " " + task.BizId)
						buttons = append(buttons, button)
					}
					keyboard := tgbotapi.NewOneTimeReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							buttons...,
						),
					)
					msg.ReplyMarkup = keyboard
				}
			case strings.Contains(update.Message.Text, model.CREATE_TASK):
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
					resp = "Please specify task ID, goal and description, in format: /edit_task <task_id> <task_goal> <task_description>"
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
					resp = "Please register first, type " + model.CREATE_USER + " <login> <password> or connect user " + model.CONNECT_USER + " <login> <password>"
				} else {
					resp = "Hello, " + update.Message.From.UserName + "!" + " You already registered as " + userID + ", to get help type " + model.HELP
				}
			case strings.Contains(update.Message.Text, model.HELP):
				resp = model.Commands[model.HELP]
			default:
				resp = "Hello, " + update.Message.From.UserName + "!" + " You said: " + update.Message.Text + ", to get help type " + model.HELP
			}

			msg.ReplyToMessageID = update.Message.MessageID
			msg.Text = resp

			//if len(resp) > 4095 {
			//	// split message and send in parts
			//	for len(resp) > 4095 {
			//		_, err = t.Send(update.Message.Chat.ID, update.Message.MessageID, resp[:4095])
			//		if err != nil {
			//			log.Err(err).Msg("send error")
			//		}
			//		resp = resp[4095:]
			//	}
			//}
			_, err = t.bot.Send(msg)
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
