package dbclient

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/google/uuid"
	"rpgMonster/internal/model"
)

const (
	LOGIN = "login"
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

func (m *Manager) CheckPassword(login string, password string) (id string, err error) {
	var user model.User
	err = m.collectionUsers.FindOne(context.TODO(), map[string]string{LOGIN: login}).Decode(&user)
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
