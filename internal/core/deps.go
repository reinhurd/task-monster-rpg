//go:generate mockgen -source=deps.go -destination=mock_test.go -package=core
package core

import (
	"context"

	"rpgMonster/internal/model"
)

type GPTClient interface {
	GetCompletion(systemContent, userContent string) (model.GPTAnswer, error)
}

type DBClient interface {
	CreateTask(ctx context.Context, task *model.Task) (err error)
	GetTask(ctx context.Context, bizID string) (task *model.Task, err error)
	GetTaskByIDAndUserID(ctx context.Context, taskID int64, userID string) (task *model.Task, err error)
	UpdateTask(ctx context.Context, task *model.Task) (err error)
	CreateNewUser(login, password string) (id string, err error)
	CreateNewUserTG(login, password string, telegramID int64) (id string, err error)
	CheckPassword(login, password string) (id string, tempToken string, err error)
	GetUserByTempToken(tempToken string) (id string, err error)
	CheckUserByLogin(login string) (id string, err error)
	CheckUserByBizID(bizID string) (id string, err error)
	GetUserByTGID(telegramID int64) (id string, err error)
	UpdateUserTGID(userID string, telegramID int64) error
	GetTaskListByUserID(userID string) (tasks []model.Task, err error)
}
