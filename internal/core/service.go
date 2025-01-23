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

func (s *Service) GetTask(ctx context.Context, bizID string, userID string) (task *model.Task, err error) {
	//todo check rights of executor or reviewer with user ID
	task, err = s.dbManager.GetTask(ctx, bizID)
	if err != nil {
		return nil, err
	}
	if task.Executor != userID || (task.Reviewer != nil && *task.Reviewer != userID) {
		return nil, fmt.Errorf("no rights to view task")
	}
	return task, nil
}

func (s *Service) GetListTasksByUserID(ctx context.Context, userID string) (tasks []model.Task, err error) {
	return s.dbManager.GetTaskListByUserID(userID)
}

func (s *Service) UpdateTask(ctx context.Context, task *model.Task) (err error) {
	return s.dbManager.UpdateTask(ctx, task)
}

func (s *Service) CreateTaskFromGPTByRequest(req string, userID string) (task *model.Task, err error) {
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
	task.Executor = userID
	err = s.dbManager.CreateTask(context.TODO(), task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) CreateUserFromTG(login, password string, TGID int) (id string, err error) {
	return s.dbManager.CreateNewUserTG(login, password, TGID)
}

func (s *Service) ConnectUserToTG(userID string, telegramID int) (err error) {
	err = s.dbManager.UpdateUserTGID(userID, telegramID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ValidateUserTG(telegramID int) (id string, err error) {
	//find user by telegram ID
	userID, err := s.dbManager.GetUserByTGID(telegramID)
	if err != nil {
		return "", err
	}
	return userID, nil
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
