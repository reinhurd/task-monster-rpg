package main

import (
	"encoding/csv"
	"github.com/gocarina/gocsv"
	"log"
	"os"
)

const PLAYERFILE = "players.csv"

// entites about gaming models of user when he got and doing tasks
// todo getters and setters for all methods
type Player struct {
	Name        string
	Token       string //must be unique
	CurrentTask string
	Level       string
	Xp          string
	Health      string //percentage
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
	cur := make([]Player, 0, 100)
	_ = gocsv.UnmarshalWithoutHeaders(f, &cur)

	return cur
}
