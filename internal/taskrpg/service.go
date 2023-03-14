package taskrpg

type Service struct {
	ios Ioservice
}

func New(ios Ioservice) *Service {
	return &Service{
		ios: ios,
	}
}
