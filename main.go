package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"log"
	"strings"
)

const (
	TOKEN = "token"
	TOPIC = "topic"
)

func main() {
	//todo add context
	//getChat("What is life?")
	initApp()
	fmt.Printf("Starting server for testing HTTP POST...\n")

	ln, err := reuseport.Listen("tcp4", "localhost:8080")
	if err != nil {
		log.Fatalf("error in reuseport listener: %v", err)
	}

	if err = fasthttp.Serve(ln, handler); err != nil {
		log.Fatalf("error in fasthttp Server: %v", err)
	}
}

func initApp() {
	generateTopics(DEFAULT_TOPICS_DATA)
	generatePlayers(DEFAULT_PLAYERS_DATA)
}

func handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	// http://localhost:8080/get_tasks?token=1&topic=php
	case "/get_tasks":
		if string(ctx.Method()) == "GET" {
			findTaskHandler(ctx)
		}
		return
	// http://localhost:8080/complete_tasks?token=1&topic=php
	case "/complete_tasks":
		if string(ctx.Method()) == "GET" {
			completeTaskHandler(ctx)
		}
		return
	default:
		ctx.Error("404 not found.", fasthttp.StatusNotFound)
		return
	}
}

func completeTaskHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	token := ctx.QueryArgs().Peek(TOKEN)
	topic := ctx.QueryArgs().Peek(TOPIC)
	player, err := validatePlayer(string(token))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	err = completeTopic(player, string(topic))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp := make(map[string]string)
	resp["result"] = fmt.Sprintf("your new level is %v your new xp is %v", player.Level, player.Xp)
	//saving
	setPlayers(player)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	fmt.Fprintf(ctx, string(jsonResp))
	return
}

func findTaskHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	token := ctx.QueryArgs().Peek(TOKEN)
	topic := ctx.QueryArgs().Peek(TOPIC)

	pl, err := validatePlayer(string(token))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	if len(topic) < 1 {
		ctx.Error("400 bad request", fasthttp.StatusBadRequest)
		return
	}
	resp := make(map[string]string)
	//set a random topic to a player
	curRandTopic, err := findTopic(string(token), string(topic))
	if err != nil {
		ctx.Error("400 bad request", fasthttp.StatusBadRequest)
		log.Fatalf("Err: %s", err)
		return
	}
	resp["result"] = fmt.Sprintf("ok for token: %v, topic: %v == %v", token, topic, curRandTopic)
	setTopicAndRemoveOldToPlayer(curRandTopic, pl)

	setPlayers(pl)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	fmt.Fprintf(ctx, string(jsonResp))
	return
}

func validatePlayer(token string) (*Player, error) {
	if len(token) < 1 {
		return nil, errors.New("no token in input")
	}
	players := loadPlayers()
	for _, player := range players {
		if strings.ToLower(player.Token) == strings.ToLower(token) {
			return &player, nil
		}
	}
	return nil, errors.New("no token found for players")
}
