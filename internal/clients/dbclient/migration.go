package dbclient

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"rpgMonster/internal/model"
)

func (m *Manager) CreateCollections() {
	err := m.db.CreateCollection(context.Background(), model.TASKS_COLLECTION)
	if err != nil {
		panic(err)
	}
	err = m.db.CreateCollection(context.Background(), model.USERS_COLLECTION)
	if err != nil {
		panic(err)
	}
}

func (m *Manager) DebugMigration() error {
	//find if user exists
	var testUserID string
	user := &model.User{Login: "test"}
	err := m.collectionUsers.FindOne(context.Background(), nil).Decode(user)
	if err != nil {
		log.Info().Msg("User not found, creating new user")
		testUserID, err = m.CreateNewUser("test", "test")
		if err != nil {
			return err
		}
	} else {
		testUserID = user.BizID
	}

	//count task for test user, create one if none found
	filter := bson.M{"executor": testUserID}
	count, err := m.collectionTasks.CountDocuments(context.Background(), filter)
	if err != nil {
		return err
	}
	if count == 0 {
		testTask := &model.Task{
			Title:       "Test Title",
			Description: "Test description",
			Executor:    testUserID,
			Reviewer:    &testUserID,
			Completed:   false,
			Deadline:    time.Now().Add(time.Hour * 24),
			Tags:        []string{"test", "task"},
		}
		err = m.CreateTask(context.Background(), testTask)
		if err != nil {
			return err
		}
	}
	return nil
}
