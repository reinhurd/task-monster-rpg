package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	TOKEN = "token"
	TOPIC = "topic"
)

func main() {
	//todo add context
	//getChat()
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

func getChat() {
	client := &http.Client{}
	var data = strings.NewReader(`{
  "model": "text-davinci-003",
  "prompt": "What is php?",      
  "max_tokens": 4000,
  "temperature": 1.0
}`)
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/completions", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik1UaEVOVUpHTkVNMVFURTRNMEZCTWpkQ05UZzVNRFUxUlRVd1FVSkRNRU13UmtGRVFrRXpSZyJ9.eyJodHRwczovL2FwaS5vcGVuYWkuY29tL3Byb2ZpbGUiOnsiZW1haWwiOiJyb2FydXNAeWFuZGV4LnJ1IiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImdlb2lwX2NvdW50cnkiOiJBUiJ9LCJodHRwczovL2FwaS5vcGVuYWkuY29tL2F1dGgiOnsidXNlcl9pZCI6InVzZXItNm5SWWhRazdWWVZ4d1hRTGVjUDM5NzlrIn0sImlzcyI6Imh0dHBzOi8vYXV0aDAub3BlbmFpLmNvbS8iLCJzdWIiOiJhdXRoMHw2M2MzMzEwMTQyOWQxM2ZhOTg5MGRkNGIiLCJhdWQiOlsiaHR0cHM6Ly9hcGkub3BlbmFpLmNvbS92MSIsImh0dHBzOi8vb3BlbmFpLmF1dGgwLmNvbS91c2VyaW5mbyJdLCJpYXQiOjE2NzQ3NjgyNDksImV4cCI6MTY3NTM3MzA0OSwiYXpwIjoiVGRKSWNiZTE2V29USHROOTVueXl3aDVFNHlPbzZJdEciLCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIG1vZGVsLnJlYWQgbW9kZWwucmVxdWVzdCBvcmdhbml6YXRpb24ucmVhZCBvZmZsaW5lX2FjY2VzcyJ9.LUcn9JjP7DWa9oPKHcKb0jXKWcQrcm3V5kMGEch4na8Y8GiScri3uJZuVGPOf0APHqPGXMt3-dKVWylNj8C7TcJjyjPkACp-9nv1UACbQ2j0ORN2cCXhfNmzmCOCWxxjZ2ACPagtblMRZrybxv8k3X7BU9eckGVVeWFpKhenihaNPrN4slusGMaqgX2b7z1NGUZC4MOHKTQqvsjAIXSERsDlvsJXO8BbS3G0PuDxqyookgd4ca30QaWf4xoEVIBoUpWyGEFfDtVwW18bByMICjPZLvHoxTIqCz92UeGnzsH2lZn7x86h7O06WHw85aRu9etqlAj8FNtRfbk5C5rj0w")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
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
