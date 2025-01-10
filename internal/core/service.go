package core

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"rpgMonster/internal/model"
)

type Service struct {
	gptClient GPTClient
	dbManager DBClient
}

func (s *Service) DoSomething() string {
	return "Hello, world!"
}

func (s *Service) CreateTask(ctx context.Context, task *model.Task) (err error) {
	return s.dbManager.CreateTask(ctx, task)
}

func (s *Service) UpdateTask(ctx context.Context, task *model.Task) (err error) {
	return s.dbManager.UpdateTask(ctx, task)
}

func (s *Service) CreateTaskFromGPTByRequest(req string) (task *model.Task, err error) {
	if req == "" {
		return nil, fmt.Errorf("empty request")
	}
	//set goal
	goal := "learn " + req
	resp, err := s.gptClient.GetCompletion(model.GPT_SYSTEM_PROMPT, fmt.Sprintf(model.GPT_DEFAULT_REQUEST, goal))
	if err != nil {
		log.Error().Err(err).Msg("error getting completion")
		return
	}
	//todo add user ID somehow
	task = &model.Task{}
	task.Title = goal
	task.Description = resp.Choices[0].Message.Content
	err = s.dbManager.CreateTask(context.TODO(), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) CreateNewUser(login, password string) (id string, err error) {
	return s.dbManager.CreateNewUser(login, password)
}

func (s *Service) CheckPassword(login, password string) (id string, err error) {
	return s.dbManager.CheckPassword(login, password)
}

func NewService(gptClient GPTClient, taskManager DBClient) *Service {
	return &Service{
		gptClient: gptClient,
		dbManager: taskManager,
	}
}
