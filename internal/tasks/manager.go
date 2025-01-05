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
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	BizId       string             `bson:"biz_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Executor    string             `bson:"executor"` // ID of the user executing the task
	Reviewer    *string            `bson:"reviewer"` // Optional ID of the reviewing user
	Completed   bool               `bson:"completed"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type Manager struct {
	collection *mongo.Collection
}

func NewManager() *Manager {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Err(err).Msg("Failed to connect to MongoDB")
	}

	db := client.Database("checklist")
	return &Manager{
		collection: db.Collection("tasks"),
	}
}

func CreateTask(ctx context.Context, task *Task) error {
	m := NewManager()
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

func GetTask(ctx context.Context, id primitive.ObjectID) (*Task, error) {
	var task Task
	m := NewManager()
	err := m.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func UpdateTask(ctx context.Context, task *Task) error {
	task.UpdatedAt = time.Now()
	m := NewManager()
	_, err := m.collection.UpdateOne(
		ctx,
		bson.M{"biz_id": task.BizId},
		bson.M{"$set": bson.M{"completed": task.Completed}},
	)
	return err
}

func DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	m := NewManager()
	_, err := m.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func ListTasks(ctx context.Context) ([]Task, error) {
	m := NewManager()
	cursor, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Get tasks by executor
func GetTasksByExecutor(ctx context.Context, executorID string) ([]Task, error) {
	m := NewManager()
	cursor, err := m.collection.Find(ctx, bson.M{"executor": executorID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Get tasks by reviewer
func GetTasksByReviewer(ctx context.Context, reviewerID string) ([]Task, error) {
	m := NewManager()
	cursor, err := m.collection.Find(ctx, bson.M{"reviewer": reviewerID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}
