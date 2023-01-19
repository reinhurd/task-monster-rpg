package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
)

const (
	TOKEN = "token"
	TOPIC = "topic"
)

// http://localhost:8080/get_tasks?token=1&topic=php
func main() {
	// unite with crete init methods
	generateTopics()
	generatePlayers()
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
	if string(ctx.Path()) == "/get_tasks" && string(ctx.Method()) == "GET" {
		findTaskHandler(ctx)
		return
	}
	if string(ctx.Path()) == "/complete_tasks" && string(ctx.Method()) == "GET" {
		//todo create some architecture for completing tasks
		return
	}
	ctx.Error("404 not found.", fasthttp.StatusNotFound)
	return
}

func findTaskHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	//todo validate
	token := ctx.QueryArgs().Peek(TOKEN)
	topic := ctx.QueryArgs().Peek(TOPIC)
	err := validate(string(token))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	if len(topic) < 1 {
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
	return
}

func validate(token string) error {
	if len(token) < 1 {
		return errors.New("no token in input")
	}
	players := loadPlayers()
	for _, player := range players {
		// normalize
		if player.Token == token {
			return nil
		}
	}
	return errors.New("no token found for players")
}
