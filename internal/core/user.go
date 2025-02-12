package core

func (s *Service) CreateUserFromTG(login, password string, TGID int64) (id string, err error) {
	return s.dbManager.CreateNewUserTG(login, password, TGID)
}

func (s *Service) ConnectUserToTG(userID string, telegramID int64) (err error) {
	err = s.dbManager.UpdateUserTGID(userID, telegramID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ValidateUserTG(telegramID int64) (id string, err error) {
	//find user by telegram ID
	userID, err := s.dbManager.GetUserByTGID(telegramID)
	if err != nil {
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
