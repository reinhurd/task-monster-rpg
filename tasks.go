package main

import (
	"math/rand"
	"strings"
)

func findTopic(token, topic string) string {
	topics := getTopics()
	//todo normalize searched values
	if val, ok := topics[topic]; ok {
		return "ok for " + token + " topic " + topic + " " + findRandomTasks(val)
	}
	return topic + " not found"
}

func findRandomTasks(tasks string) string {
	splStrings := strings.Split(tasks, ",")
	if len(splStrings) == 1 {
		return splStrings[0]
	}
	return splStrings[rand.Intn(len(splStrings))]
}
