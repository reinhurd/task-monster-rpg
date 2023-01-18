package main

import (
	"github.com/gocarina/gocsv"

	"encoding/csv"
	"log"
	"os"
)

const TOPICFILE = "topics.csv"

type Topic struct {
	MainTheme string
	Topics    string
}

func generateTopics() {
	topics := [][]string{
		{"Main Theme", "Topic"},
		{"golang", "Concurrency,Parallelism,Goroutine,Frameworks"},
		{"php", "Concurrency,Parallelism,PHP9,Frameworks"},
	}

	csvFile, err := os.Create(TOPICFILE)

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

func getTopics() map[string]string {
	f, err := os.Open(TOPICFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	cur := make([]Topic, 0, 100)
	_ = gocsv.UnmarshalWithoutHeaders(f, &cur)
	res := make(map[string]string, len(cur))
	for i := range cur {
		if i == 0 {
			continue
		}
		res[cur[i].MainTheme] = cur[i].Topics
	}

	return res
}
