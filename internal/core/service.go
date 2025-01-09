package core

import (
	"context"

	"github.com/rs/zerolog/log"
	"rpgMonster/internal/clients/gpt"
	"rpgMonster/internal/clients/telegram"
	"rpgMonster/internal/model"
	"rpgMonster/internal/tasks"
)

const (
	systemPrompt = "You a personal assistant, helping people to set concrete detailed steps to achieve goals"
)

type Service struct {
	gptClient   *gpt.Client
	taskManager *tasks.Manager
	tgBot       *telegram.TGBot
}

func (s *Service) DoSomething() string {
	return "Hello, world!"
}

func (s *Service) CreateTask(ctx context.Context, task *model.Task) (err error) {
	return s.taskManager.CreateTask(ctx, task)
}

func (s *Service) UpdateTask(ctx context.Context, task *model.Task) (err error) {
	return s.taskManager.UpdateTask(ctx, task)
}

func (s *Service) RunTG() {
	updChan := s.tgBot.GetUpdatesChan()
	err := s.tgBot.HandleUpdate(updChan)
	if err != nil {
		panic(err)
	}
}

func (s *Service) CreateTaskFromGPTByRequest(req string) (task *model.Task, err error) {
	//set goal
	goal := "learn " + req
	resp, err := s.gptClient.GetCompletion(systemPrompt, "Write a one single daily task to achieve goal "+goal+
		", in format: 'daily task: task description: requirements to check' and delimiter is comma")
	if err != nil {
		log.Error().Err(err).Msg("error getting completion")
		return
	}
	//todo add user ID somehow
	task = &model.Task{}
	task.Title = goal
	task.Description = resp.Choices[0].Message.Content
	err = s.taskManager.CreateTask(context.TODO(), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func NewService(gptClient *gpt.Client, taskManager *tasks.Manager, tgBot *telegram.TGBot) *Service {
	return &Service{
		gptClient:   gptClient,
		taskManager: taskManager,
		tgBot:       tgBot,
	}
}
