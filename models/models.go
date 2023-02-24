package models

// entites about gaming models of user when he got and doing tasks
type PlayerDTO struct {
	Name        string
	Token       string //must be unique
	CurrentTask string
	Level       string
	Xp          string
	Health      string //percentage
}
