package tasks

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rpgMonster/internal/model"
)

const (
	BIZ_ID    = "biz_id"
	COMPLETED = "completed"
	EXECUTOR  = "executor"
	REVIEWER  = "reviewer"
)

type Manager struct {
	collection *mongo.Collection
}

// todo move to db_client pkg
func NewManager() *Manager {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Err(err).Msg("Failed to connect to MongoDB")
	}

	db := client.Database(model.DB_NAME)
	return &Manager{
		collection: db.Collection(model.TASKS_COLLECTION),
	}
}

func (m *Manager) CreateTask(ctx context.Context, task *model.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.BizId = uuid.New().String()

	result, err := m.collection.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	task.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (m *Manager) GetTask(ctx context.Context, bizID string) (task *model.Task, err error) {
	err = m.collection.FindOne(ctx, bson.M{BIZ_ID: bizID}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (m *Manager) UpdateTask(ctx context.Context, task *model.Task) error {
	task.UpdatedAt = time.Now()
	_, err := m.collection.UpdateOne(
		ctx,
		bson.M{BIZ_ID: task.BizId},
		bson.M{"$set": bson.M{COMPLETED: task.Completed}},
	)
	return err
}

func (m *Manager) DeleteTask(ctx context.Context, bizID string) error {
	_, err := m.collection.DeleteOne(ctx, bson.M{BIZ_ID: bizID})
	return err
}

func (m *Manager) ListTasks(ctx context.Context) (tasks []model.Task, err error) {
	cursor, err := m.collection.Find(ctx, bson.M{})
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
	cursor, err := m.collection.Find(ctx, bson.M{EXECUTOR: executorID})
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
	cursor, err := m.collection.Find(ctx, bson.M{REVIEWER: reviewerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
