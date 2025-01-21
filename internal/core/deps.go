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
	UpdateTask(ctx context.Context, task *model.Task) (err error)
	CreateNewUser(login, password string) (id string, err error)
	CreateNewUserTG(login, password string, TGID int) (id string, err error)
	CheckPassword(login, password string) (id string, err error)
	GetUserByTGID(telegramID int) (id string, err error)
	UpdateUserTGID(userID string, telegramID int) error
}
