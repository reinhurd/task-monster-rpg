package main

import (
	"rpgMonster/internal/ioservice"
	"strings"
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
	ios := ioservice.New()
	req := make([][]string, 0)
	for _, topic := range topics {
		req = append(req, topic.ToCSV())
	}
	ios.SaveTopics(TOPICFILE, req)
}

func getTopics() []Topic {
	ios := ioservice.New()
	resRaw := ios.GetTopics(TOPICFILE)
	res := make([]Topic, 0)
	for _, t := range resRaw {
		res = append(res, Topic{
			MainTheme: t.MainTheme,
			Topics:    t.Topics,
		})
	}

	return res
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

func saveNewTopics(theme string, topics string) error {
	allTopics := getTopics()
	for _, t := range allTopics {
		if strings.ToLower(t.MainTheme) == strings.ToLower(theme) {
			t.Topics = strings.ToLower(topics)
			saveTopics(allTopics)

			return nil
		}
	}

	newTopic := Topic{
		MainTheme: theme,
		Topics:    topics,
	}

	topicsToSave := append(allTopics, newTopic)
	saveTopics(topicsToSave)

	return nil
}
