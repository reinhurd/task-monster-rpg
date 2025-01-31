package core

type Service struct {
	gptClient GPTClient
	dbManager DBClient
}

func (s *Service) DoSomething() string {
	return "Hello, world!"
}

func NewService(gptClient GPTClient, taskManager DBClient) *Service {
	return &Service{
		gptClient: gptClient,
		dbManager: taskManager,
	}
}
