package main

import (
	"github.com/gocarina/gocsv"

	"encoding/csv"
	"log"
	"os"
)

const TOPICFILE = "topics.csv"

var DEFAULT_TOPICS = []Topic{{
	MainTheme: "Main Theme",
	Topics:    "Topic",
}, {
	MainTheme: "golang",
	Topics:    "Concurrency,Parallelism,Goroutine,Frameworks",
}, {
	MainTheme: "php",
	Topics:    "Concurrency,Parallelism,PHP9,Frameworks",
}}

type Topic struct {
	MainTheme string
	Topics    string
}

func (t *Topic) ToCSV() []string {
	return []string{t.MainTheme, t.Topics}
}

func saveTopics(topics []Topic) {
	csvFile, err := os.Create(TOPICFILE)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, topic := range topics {
		_ = csvwriter.Write(topic.ToCSV())
	}
	csvwriter.Flush()
	csvFile.Close()
}

func getTopics() []Topic {
	f, err := os.Open(TOPICFILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	cur := make([]Topic, 0, 100)
	_ = gocsv.UnmarshalWithoutHeaders(f, &cur)

	return cur
}

func makeTopicsAsMap(cur []Topic) map[string]string {
	res := make(map[string]string, len(cur))
	for i := range cur {
		if i == 0 {
			continue
		}
		res[cur[i].MainTheme] = cur[i].Topics
	}

	return res
}
