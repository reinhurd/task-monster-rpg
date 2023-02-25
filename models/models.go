package models

type PlayerDTO struct {
	Name        string
	Token       string //must be unique
	CurrentTask string
	Level       string
	Xp          string
	Health      string //percentage
}

type TopicDTO struct {
	MainTheme string
	Topics    string
}
