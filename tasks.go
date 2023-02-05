package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func findTopic(topic string) (string, error) {
	topics := makeTopicsAsMap(getTopics())
	if val, ok := topics[topic]; ok {
		return findRandomTasks(val), nil
	}
	return "", fmt.Errorf("%v not found", topic)
}

func findRandomTasks(tasks string) string {
	splStrings := strings.Split(tasks, ",")
	if len(splStrings) == 1 {
		return splStrings[0]
	}
	return splStrings[rand.Intn(len(splStrings)-1)]
}
