package main

import (
	"encoding/csv"
	"errors"
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"strconv"
)

const PLAYERFILE = "players.csv"
const DEFAULT_REWARD = 10

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

func setNewLevel(level, xp int64) (int64, int64) {
	newLevelXp := level * level * 1000
	if xp >= newLevelXp {
		return level + 1, xp - newLevelXp
	}
	return level, xp
}

func completeTasksForXp(xp, reward int64) int64 {
	return xp + reward
}

func completeTopic(pl *Player, topic string) error {
	//todo normalize
	if pl.CurrentTask != topic {
		return errors.New("topic is not set in player")
	}
	pl.CurrentTask = ""
	//todo test for all math cases
	pl.Xp = completeTasksForXp(pl.Xp, DEFAULT_REWARD)
	pl.Level, pl.Xp = setNewLevel(pl.Level, pl.Xp)
	return nil
}

// create db with all existed players for future use
func generatePlayers() {
	players := [][]string{
		{"name", "token", "currentTask", "level", "xp", "health"},
		{"PersonOne", "123456", "PHP", "1", "100", "100"},
		{"PersonTwo", "221459", "Golang", "1", "99", "100"},
	}

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
