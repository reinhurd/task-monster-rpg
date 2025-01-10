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

func NewManager() *Manager {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Err(err).Msg("Failed to connect to MongoDB")
	}

	db := client.Database(model.DB_NAME)
	return &Manager{
		collectionTasks: db.Collection(model.TASKS_COLLECTION),
		collectionUsers: db.Collection(model.USERS_COLLECTION),
	}
}
