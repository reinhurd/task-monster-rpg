package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func findTopic(token, topic string) (string, error) {
	topics := getTopics()
	//todo normalize searched values
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
	return splStrings[rand.Intn(len(splStrings))]
}
