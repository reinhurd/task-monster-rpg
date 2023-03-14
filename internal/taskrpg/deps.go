package taskrpg

import "rpgMonster/models"

type Ioservice interface {
	SavePlayers(file string, players [][]string)
	LoadPlayers(file string) []models.PlayerDTO
	SaveTopics(file string, topics [][]string)
	GetTopics(file string) []models.TopicDTO
}
