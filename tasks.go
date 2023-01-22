package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func findTopic(token, topic string) string {
	topics := getTopics()
	//todo normalize searched values
	if val, ok := topics[topic]; ok {
		return fmt.Sprintf("ok for token: %v, topic: %v == %v", token, topic, findRandomTasks(val))
	}
	return fmt.Sprintf("%v not found", topic)
}

func findRandomTasks(tasks string) string {
	splStrings := strings.Split(tasks, ",")
	if len(splStrings) == 1 {
		return splStrings[0]
	}
	return splStrings[rand.Intn(len(splStrings))]
}
