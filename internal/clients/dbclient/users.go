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
	LOGIN = "login"
	TGID  = "telegram_id"
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

func (m *Manager) CreateNewUserTG(login, password string, telegramID int) (id string, err error) {
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

// todo maybe set temptoken?
func (m *Manager) CheckPassword(login string, password string) (id string, err error) {
	var user model.User
	err = m.collectionUsers.FindOne(context.TODO(), bson.M{LOGIN: login}).Decode(&user)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	hash.Write(salt)
	hash.Write([]byte(password))
	if user.Password != hex.EncodeToString(hash.Sum(nil)) {
		return "", fmt.Errorf("invalid password")
	}

	return user.BizID, nil
}
