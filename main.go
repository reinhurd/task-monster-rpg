package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
	"math/rand"
	"strings"
)

const (
	TOKEN = "token"
	TOPIC = "topic"
)

// http://localhost:8080/get_tasks?token=1&topic=php
func main() {
	generateTopics()
	fmt.Printf("Starting server for testing HTTP POST...\n")

	ln, err := reuseport.Listen("tcp4", "localhost:8080")
	if err != nil {
		log.Fatalf("error in reuseport listener: %v", err)
	}

	if err = fasthttp.Serve(ln, handler); err != nil {
		log.Fatalf("error in fasthttp Server: %v", err)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	if string(ctx.Path()) != "/get_tasks" || string(ctx.Method()) != "GET" {
		ctx.Error("404 not found.", fasthttp.StatusNotFound)
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	//todo validate
	token := ctx.QueryArgs().Peek(TOKEN)
	topic := ctx.QueryArgs().Peek(TOPIC)
	if len(token) < 1 || len(topic) < 1 {
		ctx.Error("400 bad request", fasthttp.StatusBadRequest)
		return
	}
	resp := make(map[string]string)
	resp["result"] = findTopic(string(token), string(topic))
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	fmt.Fprintf(ctx, string(jsonResp))
}

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
