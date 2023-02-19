package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"strconv"
	"strings"
)

const PLAYERFILE = "players.csv"
const DEFAULT_REWARD = 10
const DEFAULT_FINE = 20
const DEFAULT_TOKEN_LENGHT = 10

var DEFAULT_PLAYERS_DATA = [][]string{
	{"name", "token", "currentTask", "level", "xp", "health"},
	{"PersonOne", "123456", "PHP", "1", "100", "100"},
	{"PersonTwo", "221459", "Golang", "1", "99", "100"},
}

// entites about gaming models of user when he got and doing tasks
type PlayerDTO struct {
	Name        string
	Token       string //must be unique
	CurrentTask string
	Level       string
	Xp          string
	Health      string //percentage
}

type Player struct {
	Name        string
	Token       string //must be unique
	CurrentTask string
	Level       int64
	Xp          int64
	Health      int64 //percentage
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

func (p *Player) completeTopic(topic string) error {
	if strings.ToLower(p.CurrentTask) != strings.ToLower(topic) {
		return errors.New("topic is not set in player")
	}
	p.CurrentTask = ""
	//todo test for all math cases
	p.completeTasksForXp(DEFAULT_REWARD)
	p.setNewLevel()
	return nil
}

// create db with all existed players for future use
func generatePlayers(players [][]string) {
	csvFile, err := os.Create(PLAYERFILE)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, player := range players {
		_ = csvwriter.Write(player)
	}
	csvwriter.Flush()
	csvFile.Close()
}

func loadPlayers() []Player {
	f, err := os.Open(PLAYERFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	cur := make([]PlayerDTO, 0, 100)
	_ = gocsv.UnmarshalWithoutHeaders(f, &cur)

	return toPlayer(cur)
}

func validatePlayerName(name string) error {
	pls := loadPlayers()
	for _, oldPl := range pls {
		if strings.ToLower(oldPl.Name) == strings.ToLower(name) {
			return fmt.Errorf("player name %s already exists", name)
		}
	}
	return nil
}

func createNewPlayer(name string) Player {
	player := Player{
		Name:        name,
		Token:       generateToken(),
		CurrentTask: "",
		Level:       1,
		Xp:          0,
		Health:      0,
	}
	setPlayers(&player)
	return player
}

func generateToken() string {
	return RandStringBytesMaskImpr(DEFAULT_TOKEN_LENGHT)
}

func setPlayers(plr *Player) {
	players := loadPlayers()
	resPlrs := make([]Player, 0, len(players))
	for _, oldPl := range players {
		//todo test this
		if oldPl.Token != plr.Token {
			resPlrs = append(resPlrs, oldPl)
		}
	}
	resPlrs = append(resPlrs, *plr)

	csvFile, err := os.Create(PLAYERFILE)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, player := range resPlrs {
		_ = csvwriter.Write(player.toCSV())
	}
	csvwriter.Flush()
	csvFile.Close()
}

func setTopicAndRemoveOldToPlayer(topic string, pl *Player) {
	if topic == "" {
		return
	}
	if pl.CurrentTask != "" {
		//fine player for undoing task
		pl.Xp = pl.Xp - DEFAULT_FINE
		fmt.Printf("The player %s was fined by amount %v for not completed task", pl.Name, DEFAULT_FINE)
	}
	pl.CurrentTask = strings.ToLower(topic)
}

func toPlayer(plDto []PlayerDTO) []Player {
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
