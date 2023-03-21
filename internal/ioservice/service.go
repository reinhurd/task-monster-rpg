package ioservice

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocarina/gocsv"

	"rpgMonster/models"
)

type service struct{}

func New() *service {
	return &service{}
}

// SavePlayers create db with all existed players for future use
func (s *service) SavePlayers(file string, players [][]string) {
	csvFile, err := os.Create(file)

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

func (s *service) LoadPlayers(file string) []models.PlayerDTO {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	cur := make([]models.PlayerDTO, 0, 100)
	_ = gocsv.UnmarshalWithoutHeaders(f, &cur)

	return cur
}

func (s *service) SaveTopics(file string, topics [][]string) {
	csvFile, err := os.Create(file)

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

func (s *service) GetTopics(file string) []models.TopicDTO {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	cur := make([]models.TopicDTO, 0, 100)
	_ = gocsv.UnmarshalWithoutHeaders(f, &cur)

	return cur
}
