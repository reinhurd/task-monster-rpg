package dbclient

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"rpgMonster/internal/model"
)

const (
	LOGIN      = "login"
	TGID       = "telegram_id"
	TEMP_TOKEN = "temp_token"
)

func (m *Manager) CreateNewUser(login string, password string) (id string, err error) {
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(password))
	passwordHash := hex.EncodeToString(hash.Sum(nil))
	saltHex := hex.EncodeToString(salt)

	var user = model.User{
		BizID:    uuid.New().String(),
		Login:    login,
		Password: passwordHash,
		Salt:     saltHex,
	}

	_, err = m.collectionUsers.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return user.BizID, nil
}

func (m *Manager) CreateNewUserTG(login, password string, telegramID int64) (id string, err error) {
	userID, err := m.CreateNewUser(login, password)
	if err != nil {
		return "", err
	}
	//update user with TGID
	_, err = m.collectionUsers.UpdateOne(context.TODO(), bson.M{BIZ_ID: userID}, bson.M{"$set": bson.M{TGID: telegramID}})
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (m *Manager) GetUserByTGID(telegramID int64) (id string, err error) {
	var user model.User
	err = m.collectionUsers.FindOne(context.TODO(), bson.M{TGID: telegramID}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.BizID, nil
}

func (m *Manager) UpdateUserTGID(userID string, telegramID int64) error {
	_, err := m.collectionUsers.UpdateOne(context.TODO(), bson.M{BIZ_ID: userID}, bson.M{"$set": bson.M{TGID: telegramID}})
	return err
}

func (m *Manager) GetTaskListByUserID(userID string) (tasks []model.Task, err error) {
	cursor, err := m.collectionTasks.Find(context.TODO(), bson.M{EXECUTOR: userID})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.TODO(), &tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (m *Manager) CheckPassword(login string, password string) (id string, tempToken string, err error) {
	var user model.User
	err = m.collectionUsers.FindOne(context.TODO(), bson.M{LOGIN: login}).Decode(&user)
	if err != nil {
		return "", "", err
	}

	hash := sha256.New()
	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		log.Fatal(err)
		return "", "", err
	}

	hash.Write(salt)
	hash.Write([]byte(password))
	if user.Password != hex.EncodeToString(hash.Sum(nil)) {
		return "", "", fmt.Errorf("invalid password")
	}
	tempToken = getTempToken()
	//set temp token to user
	_, err = m.collectionUsers.UpdateOne(context.TODO(), bson.M{LOGIN: login}, bson.M{"$set": bson.M{TEMP_TOKEN: tempToken}})
	if err != nil {
		return "", "", err
	}

	return user.BizID, tempToken, nil
}

func (m *Manager) GetUserByTempToken(tempToken string) (id string, err error) {
	var user model.User
	err = m.collectionUsers.FindOne(context.TODO(), bson.M{TEMP_TOKEN: tempToken}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.BizID, nil
}

func (m *Manager) CheckUserByLogin(login string) (id string, err error) {
	var user model.User
	err = m.collectionUsers.FindOne(context.TODO(), bson.M{LOGIN: login}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.BizID, nil
}

func (m *Manager) CheckUserByBizID(bizID string) (id string, err error) {
	var user model.User
	err = m.collectionUsers.FindOne(context.TODO(), bson.M{BIZ_ID: bizID}).Decode(&user)
	if err != nil {
		return "", err
	}
	return user.BizID, nil
}

func getTempToken() string {
	return uuid.New().String()
}
