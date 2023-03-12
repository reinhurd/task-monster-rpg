package taskrpg

import (
	"fmt"
	"math/rand"
	"strings"
)

func (s *Service) FindTopic(topic string) (string, error) {
	topics := s.makeTopicsAsMap(s.getTopics())
	if val, ok := topics[topic]; ok {
		return s.findRandomTasks(val), nil
	}
	return "", fmt.Errorf("%v not found", topic)
}

func (s *Service) findRandomTasks(tasks string) string {
	splStrings := strings.Split(tasks, ",")
	if len(splStrings) == 1 {
		return splStrings[0]
	}
	return splStrings[rand.Intn(len(splStrings)-1)]
}
