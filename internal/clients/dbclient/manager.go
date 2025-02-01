package dbclient

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rpgMonster/internal/model"
)

type Manager struct {
	collectionTasks *mongo.Collection
	collectionUsers *mongo.Collection
}

// ping db
func (m *Manager) Ping() error {
	err := m.collectionTasks.Database().Client().Ping(context.Background(), nil)
	if err != nil {
		log.Err(err).Msg("Failed to ping MongoDB")
		return err
	}
	err = m.collectionUsers.Database().Client().Ping(context.Background(), nil)
	if err != nil {
		log.Err(err).Msg("Failed to ping MongoDB")
	}
	return err
}

func NewManager() *Manager {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Err(err).Msg("Failed to connect to MongoDB")
	}

	db := client.Database(model.DB_NAME)
	err = db.CreateCollection(context.Background(), model.TASKS_COLLECTION)
	if err != nil {
		panic(err)
	}
	err = db.CreateCollection(context.Background(), model.USERS_COLLECTION)
	if err != nil {
		panic(err)
	}
	return &Manager{
		collectionTasks: db.Collection(model.TASKS_COLLECTION),
		collectionUsers: db.Collection(model.USERS_COLLECTION),
	}
}
