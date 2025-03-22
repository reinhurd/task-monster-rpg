package core

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Service) CreateUserFromTG(login, password string, TGID int64) (id string, err error) {
	//check if user exists
	id, err = s.dbManager.GetUserByTGID(TGID)
	if id != "" {
		return "", fmt.Errorf("Your TelegramID already linked to existed userID: %v", id)
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}
	//if this user already exists
	id, err = s.dbManager.CheckUserByLogin(login)
	if id != "" {
		return "", fmt.Errorf("User with this login %s already exists with userID: %v", login, id)
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}
	return s.dbManager.CreateNewUserTG(login, password, TGID)
}

func (s *Service) ConnectUserToTG(userID string, telegramID int64) (err error) {
	//check if user exists
	id, err := s.dbManager.GetUserByTGID(telegramID)
	if id != "" {
		return fmt.Errorf("Your TelegramID already linked to another userID: %v", id)
	}
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}
	//check if user with that id exists
	id, err = s.dbManager.CheckUserByBizID(userID)
	if id == "" || err == mongo.ErrNoDocuments {
		return fmt.Errorf("User with that ID doesn't exist")
	}
	if err != nil {
		return err
	}
	err = s.dbManager.UpdateUserTGID(userID, telegramID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ValidateUserTG(telegramID int64) (id string, err error) {
	//find user by telegram ID
	userID, err := s.dbManager.GetUserByTGID(telegramID)
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
	}
	return userID, nil
}

func (s *Service) CreateNewUser(login, password string) (id string, err error) {
	return s.dbManager.CreateNewUser(login, password)
}

func (s *Service) CheckPassword(login, password string) (id string, token string, err error) {
	return s.dbManager.CheckPassword(login, password)
}

func (s *Service) GetUserByTempToken(tempToken string) (id string, err error) {
	return s.dbManager.GetUserByTempToken(tempToken)
}
