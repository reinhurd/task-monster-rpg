package taskrpg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"rpgMonster/models"
)

const PLAYERFILE = "players.csv"
const DEFAULT_REWARD = 10
const DEFAULT_FINE = 20
const DEFAULT_TOKEN_LENGHT = 10

var PLAYERS_HEADER = []string{"name", "token", "task", "level", "xp", "health"}

var DEFAULT_PLAYERS_DATA = []Player{
	{"PersonOne", "123456", "PHP", 1, 100, 100},
	{"PersonTwo", "221459", "Golang", 1, 99, 100},
}

type Player struct {
	Name        string
	Token       string // must be unique
	CurrentTask string
	Level       int64
	Xp          int64
	Health      int64 // percentage
}

func (p *Player) toCSV() []string {
	var res []string
	res = append(res, p.Name)
	res = append(res, p.Token)
	res = append(res, p.CurrentTask)
	res = append(res, strconv.Itoa(int(p.Level)))
	res = append(res, strconv.Itoa(int(p.Xp)))
	res = append(res, strconv.Itoa(int(p.Health)))

	return res
}

func (p *Player) setNewLevel() {
	newLevelXp := p.Level * p.Level * 1000
	if p.Xp >= newLevelXp {
		p.Level = p.Level + 1
		p.Xp = p.Xp - newLevelXp
	}
}

func (p *Player) completeTasksForXp(reward int64) {
	p.Xp = p.Xp + reward
}

func (p *Player) CompleteTopic(topic string) error {
	if strings.ToLower(p.CurrentTask) != strings.ToLower(topic) {
		return errors.New("topic is not set in player")
	}
	p.CurrentTask = ""
	p.completeTasksForXp(DEFAULT_REWARD)
	p.setNewLevel()
	return nil
}

func (s *Service) loadPlayers() []Player {
	cur := s.ios.LoadPlayers(PLAYERFILE)

	return toPlayer(cur)
}

func (s *Service) SavePlayers(players []Player) {
	req := make([][]string, 0)
	req = append(req, PLAYERS_HEADER)
	for _, player := range players {
		req = append(req, player.toCSV())
	}
	s.ios.SavePlayers(PLAYERFILE, req)
}

func (s *Service) ValidatePlayerName(name string) error {
	if len(name) == 0 {
		return errors.New("player name is empty")
	}
	pls := s.loadPlayers()
	for _, oldPl := range pls {
		if strings.ToLower(oldPl.Name) == strings.ToLower(name) {
			return fmt.Errorf("player name %s already exists", name)
		}
	}
	return nil
}

func (s *Service) CreateNewPlayer(name string) Player {
	player := Player{
		Name:        name,
		Token:       s.generateToken(),
		CurrentTask: "",
		Level:       1,
		Xp:          0,
		Health:      0,
	}
	s.SetPlayers(&player)
	return player
}

func (s *Service) generateToken() string {
	return RandStringBytesMaskImpr(DEFAULT_TOKEN_LENGHT)
}

func (s *Service) SetPlayers(plr *Player) {
	players := s.loadPlayers()
	resPlrs := make([]Player, 0, len(players))
	for _, oldPl := range players {
		// todo test this
		if oldPl.Token != plr.Token && oldPl.Token != "token" {
			resPlrs = append(resPlrs, oldPl)
		}
	}
	resPlrs = append(resPlrs, *plr)

	s.SavePlayers(resPlrs)
}

func (s *Service) SetTopicAndRemoveOldToPlayer(topic string, pl *Player) {
	if topic == "" {
		return
	}
	if pl.CurrentTask != "" {
		// fine player for undoing task
		pl.Xp = pl.Xp - DEFAULT_FINE
		fmt.Printf("The player %s was fined by amount %v for not completed task", pl.Name, DEFAULT_FINE)
	}
	pl.CurrentTask = strings.ToLower(topic)
}

func toPlayer(plDto []models.PlayerDTO) []Player {
	res := make([]Player, 0, len(plDto))
	for i := range plDto {
		res = append(res, Player{
			Name:        plDto[i].Name,
			Token:       plDto[i].Token,
			CurrentTask: plDto[i].CurrentTask,
			Level:       stringToInt(plDto[i].Level),
			Xp:          stringToInt(plDto[i].Xp),
			Health:      stringToInt(plDto[i].Health),
		})
	}
	return res
}

func stringToInt(s string) int64 {
	res, err := strconv.Atoi(s)
	if err == nil {
		return int64(res)
	}
	return 0
}

func (s *Service) ValidatePlayerByToken(token string) (*Player, error) {
	if len(token) < 1 {
		return nil, errors.New("no token in input")
	}
	players := s.loadPlayers()
	for _, player := range players {
		if strings.ToLower(player.Token) == strings.ToLower(token) {
			return &player, nil
		}
	}
	return nil, errors.New("no token found for players")
}

func (s *Service) ValidateTheme(th string) error {
	if len(th) < 2 {
		return errors.New("invalid theme")
	}
	return nil
}
