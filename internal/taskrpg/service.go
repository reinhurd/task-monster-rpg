package taskrpg

type Service struct {
	ios ioservice
}

func New(ios ioservice) *Service {
	return &Service{
		ios: ios,
	}
}
