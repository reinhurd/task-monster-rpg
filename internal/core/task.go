package core

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"rpgMonster/internal/model"
)

func (s *Service) CreateTask(ctx context.Context, task *model.Task) (err error) {
	return s.dbManager.CreateTask(ctx, task)
}

func (s *Service) GetTask(ctx context.Context, bizID string, userID string) (task *model.Task, err error) {
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
