package taskrpg

import (
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

func (s *Service) SaveTopics(topics []Topic) {
	if len(topics) == 0 {
		return
	}
	req := make([][]string, 0)
	for _, topic := range topics {
		req = append(req, topic.ToCSV())
	}
	s.ios.SaveTopics(TOPICFILE, req)
}

func (s *Service) getTopics() []Topic {
	resRaw := s.ios.GetTopics(TOPICFILE)

	res := make([]Topic, 0)
	for _, t := range resRaw {
		res = append(res, Topic{
			MainTheme: t.MainTheme,
			Topics:    t.Topics,
		})
	}

	return res
}

func (s *Service) makeTopicsAsMap(cur []Topic) map[string]string {
	res := make(map[string]string, len(cur))
	for i := range cur {
		if cur[i].MainTheme == DEFAULT_TOPICS[0].MainTheme {
			continue
		}
		res[cur[i].MainTheme] = cur[i].Topics
	}

	return res
}

func (s *Service) SaveNewTopics(theme string, topics string) error {
	allTopics := s.getTopics()
	for i := range allTopics {
		if strings.ToLower(allTopics[i].MainTheme) == strings.ToLower(theme) {
			allTopics[i].MainTheme = strings.ToLower(theme)
			allTopics[i].Topics = strings.ToLower(topics)
			s.SaveTopics(allTopics)

			return nil
		}
	}

	newTopic := Topic{
		MainTheme: strings.ToLower(theme),
		Topics:    strings.ToLower(topics),
	}

	topicsToSave := append(allTopics, newTopic)
	s.SaveTopics(topicsToSave)

	return nil
}
