package main

import (
	"encoding/csv"
	"log"
	"os"
)

// entites about gaming models of user when he got and doing tasks
// todo getters and setters for all methods
type Player struct {
	name        string
	token       string //must be unique
	currentTask string
	level       int64
	xp          int64
	health      int64 //percentage
}

// create db with all existed players for future use
func generatePlayers() {
	topics := [][]string{
		{"name", "token", "currentTask", "level", "xp", "health"},
		{"PersonOne", "123456", "PHP", "1", "100", "100"},
		{"PersonTwo", "221459", "Golang", "1", "99", "100"},
	}

	csvFile, err := os.Create("players.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, topic := range topics {
		_ = csvwriter.Write(topic)
	}
	csvwriter.Flush()
	csvFile.Close()
}
