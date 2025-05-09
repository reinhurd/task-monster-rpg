package telegram

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"rpgMonster/internal/core"
	"rpgMonster/internal/model"
)

// BotStore wraps a sync.Map to hold userID → lastMessage mappings.
type BotStore struct {
	lastMessages sync.Map // map[int64]string
}

// SetLastMessage stores the latest message text for a given user ID.
func (bs *BotStore) SetLastMessage(userID int64, messageText string) {
	bs.lastMessages.Store(userID, messageText)
}

// GetLastMessage retrieves the last message for a user.
// Returns the message and true if found, or "" and false otherwise.
func (bs *BotStore) GetLastMessage(userID int64) (string, bool) {
	if v, ok := bs.lastMessages.Load(userID); ok {
		return v.(string), true
	}
	return "", false
}

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
	store := &BotStore{}
	for update := range updates {
		var resp string
		if update.Message != nil { // If we got a message
			log.Info().Msgf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			userTelegramID := update.Message.From.ID

			//check if user exists and show his current tasks
			userID, err := t.svc.ValidateUserTG(userTelegramID)
			if userID == "" {
				//message that user need to register
				resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
			} else if err != nil {
				resp = err.Error()
			}

			lastChatID = update.Message.Chat.ID
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

			// get last message from cache
			lastMessage, ok := store.GetLastMessage(userTelegramID)

			var isLastMessageIsCommand bool

			//step 1. Проходим по последнему сообщению
			if ok && lastMessage != "" {
				switch {
				case strings.Contains(lastMessage, model.CREATE_TASK_GPT):
					isLastMessageIsCommand = true
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					request := update.Message.Text
					if len(request) < 1 || strings.TrimSpace(request) == "" {
						resp = "Please specify request - you need to type REQUEST - a sentence with your goal. For example: learn php."
					} else {
						task, err := t.svc.CreateTaskFromGPTByRequest(request, userID)
						if err != nil {
							resp = err.Error()
						} else {
							resp = fmt.Sprintf(model.Commands[model.CREATE_TASK_GPT], task)
						}
						store.SetLastMessage(userTelegramID, "")
					}
				case strings.Contains(lastMessage, model.CONNECT_USER):
					isLastMessageIsCommand = true
					spStr := strings.Split(update.Message.Text, " ")
					if len(spStr) < 2 {
						resp = "Please enter user login and password - for example: test1 test1"
					} else {
						userID, _, err := t.svc.CheckPassword(spStr[0], spStr[1])
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
						store.SetLastMessage(userTelegramID, "")
					}
				case strings.Contains(lastMessage, model.CREATE_TASK):
					isLastMessageIsCommand = true
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					splStr := strings.Split(update.Message.Text, " ")
					//add to second element all other elements
					if len(splStr) > 2 {
						splStr[1] = strings.Join(splStr[1:], " ")
						splStr = splStr[:2]
					}
					if len(splStr) < 2 {
						resp = "Please specify request - write title and description of task in format: <task_title> <task_description>, example: PHP this is description"
					} else {
						var task model.Task
						task.Title = splStr[0]
						task.Description = splStr[1]
						task.Executor = userID
						err = t.svc.CreateTask(context.Background(), &task)
						if err != nil {
							resp = err.Error()
						} else {
							resp = fmt.Sprintf(model.Commands[model.CREATE_TASK], task)
						}
						store.SetLastMessage(userTelegramID, "")
					}
				case strings.Contains(lastMessage, model.UPDATE_TASK_DESC):
					isLastMessageIsCommand = true
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					splStr := strings.Split(lastMessage, " ")
					//add to second element all other elements
					if len(splStr) < 2 {
						resp = "Please specify request - write description of task in format: <task_description>, example: PHP this is description"
					} else {
						taskID := splStr[1]
						taskIDint, err := strconv.Atoi(taskID)
						if err != nil {
							resp = err.Error()
							break
						}
						task, err := t.svc.GetTaskByIDAndUserID(context.Background(), int64(taskIDint), userID)
						if err != nil {
							resp = err.Error()
							break
						}
						task.Description = update.Message.Text
						err = t.svc.UpdateTask(context.Background(), task)
						if err != nil {
							resp = err.Error()
						} else {
							resp = fmt.Sprintf(model.Commands[model.UPDATE_TASK], task)
						}
						store.SetLastMessage(userTelegramID, "")
					}
				case strings.Contains(lastMessage, model.UPDATE_TASK_TITLE):
					isLastMessageIsCommand = true
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					splStr := strings.Split(lastMessage, " ")
					//add to second element all other elements
					if len(splStr) < 2 {
						resp = "Please specify request - write title of task in format: <task_title>, example: PHP"
					} else {
						taskID := splStr[1]
						taskIDint, err := strconv.Atoi(taskID)
						if err != nil {
							resp = err.Error()
							break
						}
						task, err := t.svc.GetTaskByIDAndUserID(context.Background(), int64(taskIDint), userID)
						if err != nil {
							resp = err.Error()
							break
						}
						task.Title = update.Message.Text
						err = t.svc.UpdateTask(context.Background(), task)
						if err != nil {
							resp = err.Error()
						} else {
							resp = fmt.Sprintf(model.Commands[model.UPDATE_TASK], task)
						}
						store.SetLastMessage(userTelegramID, "")
					}
				case strings.Contains(lastMessage, model.UPDATE_TASK_DATE):
					isLastMessageIsCommand = true
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					splStr := strings.Split(lastMessage, " ")
					//add to second element all other elements
					if len(splStr) < 2 {
						resp = "Please specify request - write execution date of task in DD-MM-YYYY format, example: 01-01-2023"
					} else {
						taskID := splStr[1]
						taskIDint, err := strconv.Atoi(taskID)
						if err != nil {
							resp = err.Error()
							break
						}
						task, err := t.svc.GetTaskByIDAndUserID(context.Background(), int64(taskIDint), userID)
						if err != nil {
							resp = err.Error()
							break
						}
						newDate, err := time.Parse("02-01-2006", update.Message.Text)
						if err != nil {
							resp = "Please specify request - write execution date of task in DD-MM-YYYY format, example: 01-01-2023"
						}
						task.Deadline = newDate
						err = t.svc.UpdateTask(context.Background(), task)
						if err != nil {
							resp = err.Error()
						} else {
							resp = fmt.Sprintf(model.Commands[model.UPDATE_TASK], task)
						}
						store.SetLastMessage(userTelegramID, "")
					}
				}
			}

			if !isLastMessageIsCommand {
				switch {
				case strings.Contains(update.Message.Text, model.CREATE_TASK_GPT):
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					resp = "Please specify request - you need to type REQUEST - a sentence with your goal. For example: learn php."
					store.SetLastMessage(userTelegramID, model.CREATE_TASK_GPT)
				case strings.Contains(update.Message.Text, model.CONNECT_USER):
					resp = "Please enter user login and password - for example: test1 test1"
					store.SetLastMessage(userTelegramID, model.CONNECT_USER)
				case strings.Contains(update.Message.Text, model.TASK_LIST):
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					tasks, err := t.svc.GetListTasksByUserID(context.Background(), userID)
					if err != nil {
						resp = err.Error()
					} else {
						if len(tasks) == 0 {
							resp = "You have no tasks"
							break
						} else {
							resp = "You have " + strconv.Itoa(len(tasks)) + " tasks \n"
						}
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
								} else {
									otherTasks = append(otherTasks, task)
								}
							}
						}
						buttons := make([]tgbotapi.KeyboardButton, 0)
						//set today date to DD-MM-YYYY format
						today := time.Now().Format("02-01-2006")
						resp += "---TODAY[" + today + "]---\n\n"
						for _, task := range todayTasks {
							if task.BizId != "" {
								var taskExecutionDate string
								if task.Deadline.IsZero() {
									taskExecutionDate = "not_planned"
								} else {
									taskExecutionDate = task.Deadline.Format("02-01-2006")
								}
								titleWithExecutionDate := fmt.Sprintf("%s [%s]", task.Title, taskExecutionDate)
								taskLink := model.VIEW_TASK + "_" + strconv.FormatInt(task.UnID, 10)
								resp += fmt.Sprintf("Task %v\n %v\n \n\n", titleWithExecutionDate, taskLink)
								button := tgbotapi.NewKeyboardButton(taskLink)
								buttons = append(buttons, button)
							}
						}
						resp += "---WEEK---\n\n"
						for _, task := range weekTasks {
							if task.BizId != "" {
								var taskExecutionDate string
								if task.Deadline.IsZero() {
									taskExecutionDate = "not_planned"
								} else {
									taskExecutionDate = task.Deadline.Format("02-01-2006")
								}
								titleWithExecutionDate := fmt.Sprintf("%s [%s]", task.Title, taskExecutionDate)
								taskLink := model.VIEW_TASK + "_" + strconv.FormatInt(task.UnID, 10)
								resp += fmt.Sprintf("Task %v\n %v\n \n\n", titleWithExecutionDate, taskLink)
								button := tgbotapi.NewKeyboardButton(taskLink)
								buttons = append(buttons, button)
							}
						}
						resp += "---OTHER---\n\n"
						for _, task := range otherTasks {
							if task.BizId != "" {
								var taskExecutionDate string
								if task.Deadline.IsZero() {
									taskExecutionDate = "not_planned"
								} else {
									taskExecutionDate = task.Deadline.Format("02-01-2006")
								}
								titleWithExecutionDate := fmt.Sprintf("%s [%s]", task.Title, taskExecutionDate)
								taskLink := model.VIEW_TASK + "_" + strconv.FormatInt(task.UnID, 10)
								resp += fmt.Sprintf("Task %v\n %v\n \n\n", titleWithExecutionDate, taskLink)
								button := tgbotapi.NewKeyboardButton(taskLink)
								buttons = append(buttons, button)
							}
						}
						var rows [][]tgbotapi.KeyboardButton
						for _, button := range buttons {
							rows = append(rows, tgbotapi.NewKeyboardButtonRow(button))
						}

						kb := tgbotapi.NewReplyKeyboard(
							rows...,
						)
						kb.ResizeKeyboard = true
						kb.OneTimeKeyboard = true
						msg.ReplyMarkup = kb
					}
				case strings.Contains(update.Message.Text, model.VIEW_TASK):
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					taskID := strings.Split(update.Message.Text, "_")
					if len(taskID) < 2 {
						resp = "Please specify task ID"
					} else {
						buttons := make([]tgbotapi.KeyboardButton, 0)
						taskUnID, err := strconv.ParseInt(taskID[2], 10, 64)
						if err != nil {
							resp = "Please specify task ID in format: /view_task_<task_id>"
							break
						}
						task, err := t.svc.GetTaskByIDAndUserID(context.Background(), taskUnID, userID)
						if err != nil {
							resp = err.Error()
						} else {
							resp = fmt.Sprintf(model.Commands[model.VIEW_TASK], task.BizId, task.UnID, task.Title, task.Description, task.Completed, task.Executor, task.Reviewer, task.Deadline, task.CreatedAt, task.UpdatedAt)
							button := tgbotapi.NewKeyboardButton(model.UPDATE_TASK_TITLE + " " + strconv.FormatInt(task.UnID, 10))
							buttons = append(buttons, button)
							button = tgbotapi.NewKeyboardButton(model.UPDATE_TASK_DESC + " " + strconv.FormatInt(task.UnID, 10))
							buttons = append(buttons, button)
							button = tgbotapi.NewKeyboardButton(model.UPDATE_TASK_DATE + " " + strconv.FormatInt(task.UnID, 10))
							buttons = append(buttons, button)
						}

						var rows [][]tgbotapi.KeyboardButton
						for _, button := range buttons {
							rows = append(rows, tgbotapi.NewKeyboardButtonRow(button))
						}

						kb := tgbotapi.NewReplyKeyboard(
							rows...,
						)
						kb.ResizeKeyboard = true
						kb.OneTimeKeyboard = true
						msg.ReplyMarkup = kb
					}
					break
				case strings.Contains(update.Message.Text, model.CREATE_TASK):
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					resp = "Please specify task goal and description in format: <task_goal> <task_description>, example: PHP this is description"
					store.SetLastMessage(userTelegramID, model.CREATE_TASK)
				case strings.Contains(update.Message.Text, model.UPDATE_TASK_DATE):
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					resp = "Please specify new date in format: DD-MM-YYYY, example: 01-01-2023"
					store.SetLastMessage(userTelegramID, update.Message.Text)
				case strings.Contains(update.Message.Text, model.UPDATE_TASK_TITLE):
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					resp = "Please specify new title in format: <task_title>, example: PHP"
					store.SetLastMessage(userTelegramID, update.Message.Text)
				case strings.Contains(update.Message.Text, model.UPDATE_TASK_DESC):
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					resp = "Please specify new description in format: <task_description>, example: PHP this is description"
					store.SetLastMessage(userTelegramID, update.Message.Text)
				case strings.Contains(update.Message.Text, model.UPDATE_TASK):
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp = "Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
					splStr := strings.Split(update.Message.Text, " ")
					if len(splStr) < 5 {
						resp = "Please specify task ID, goal and description, in format: /edit_task <task_id> <task_goal> <task_description> <task deadline as DD-MM-YYYY>"
					} else {
						var task model.Task
						task.BizId = splStr[1]
						task.Title = splStr[2]
						task.Description = splStr[3]
						task.Deadline, err = time.Parse("02-01-2006", splStr[4])
						if err != nil {
							resp = "Please specify task deadline in format: DD-MM-YYYY"
							break
						}
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
					userID, err := t.svc.ValidateUserTG(userTelegramID)
					if err != nil {
						resp = err.Error()
						break
					}
					if userID == "" {
						//message that user need to register
						resp += " Please register first, type " + model.CREATE_USER + " <login> <password>"
						break
					}
				}
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
