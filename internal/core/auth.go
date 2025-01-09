package core

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"rpgMonster/internal/model"
)

const (
	LOGIN = "login"
)

// GetCurrentUserID extracts the user ID from the HTTP Authorization header.
// The header should contain "Bearer <user-id>".
func GetCurrentUserID(headers http.Header) (string, error) {
	return "test-user-id", nil

	// auth := headers.Get("Authorization")
	// if auth == "" {
	// 	return "", fmt.Errorf("missing Authorization header")
	// }

	// parts := strings.Split(auth, " ")
	// if len(parts) != 2 || parts[0] != "Bearer" {
	// 	return "", fmt.Errorf("invalid Authorization header format")
	// }

	// userID := parts[1]
	// if userID == "" {
	// 	return "", fmt.Errorf("empty user ID")
	// }

	// return userID, nil
}

// todo move to db_client pkg
func CreateNewUser(login string, password string) (id string, err error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(model.DB_NAME).Collection(model.USERS_COLLECTION)

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

	_, err = collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return user.BizID, nil
}

func CheckPassword(login string, password string) (id string, err error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)

	}
	defer client.Disconnect(context.TODO())

	collection := client.Database(model.DB_NAME).Collection(model.USERS_COLLECTION)
	var user model.User
	err = collection.FindOne(context.TODO(), map[string]string{LOGIN: login}).Decode(&user)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	hash.Write([]byte(salt))
	hash.Write([]byte(password))
	if user.Password != hex.EncodeToString(hash.Sum(nil)) {
		return "", fmt.Errorf("invalid password")
	}

	return user.BizID, nil
}
