package dbclient

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"rpgMonster/internal/model"
)

const (
	BIZ_ID    = "biz_id"
	COMPLETED = "completed"
	EXECUTOR  = "executor"
	REVIEWER  = "reviewer"
)

func (m *Manager) CreateTask(ctx context.Context, task *model.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.BizId = uuid.New().String()

	result, err := m.collectionTasks.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	task.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (m *Manager) GetTask(ctx context.Context, bizID string) (task *model.Task, err error) {
	err = m.collectionTasks.FindOne(ctx, bson.M{BIZ_ID: bizID}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *Manager) UpdateTask(ctx context.Context, task *model.Task) error {
	task.UpdatedAt = time.Now()
	fmt.Println("task.BizId", task.BizId)
	_, err := m.collectionTasks.UpdateOne(
		ctx,
		bson.M{BIZ_ID: task.BizId},
		//set all fields
		bson.M{"$set": bson.M{COMPLETED: task.Completed, EXECUTOR: task.Executor, REVIEWER: task.Reviewer, "title": task.Title, "description": task.Description, "deadline": task.Deadline, "tags": task.Tags}},
	)
	return err
}

func (m *Manager) DeleteTask(ctx context.Context, bizID string) error {
	_, err := m.collectionTasks.DeleteOne(ctx, bson.M{BIZ_ID: bizID})
	return err
}

func (m *Manager) ListTasks(ctx context.Context) (tasks []model.Task, err error) {
	cursor, err := m.collectionTasks.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Get tasks by executor
func (m *Manager) GetTasksByExecutor(ctx context.Context, executorID string) (tasks []model.Task, err error) {
	cursor, err := m.collectionTasks.Find(ctx, bson.M{EXECUTOR: executorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Get tasks by reviewer
func (m *Manager) GetTasksByReviewer(ctx context.Context, reviewerID string) (tasks []model.Task, err error) {
	cursor, err := m.collectionTasks.Find(ctx, bson.M{REVIEWER: reviewerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
