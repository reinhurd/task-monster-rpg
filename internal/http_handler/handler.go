package http_handler

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/valyala/fasthttp"

	"rpgMonster/internal/chatgpt"
	"rpgMonster/internal/ioservice"
	"rpgMonster/internal/taskrpg"
)

const (
	TOKEN = "token"
	TOPIC = "topic"
	THEME = "theme"
	NAME  = "name"
)

func Handler(ctx *fasthttp.RequestCtx) {
	ios := ioservice.New()
	s := taskrpg.New(ios)

	switch string(ctx.Path()) {
	// http://localhost:8080/get_tasks?token=1&topic=php
	case "/get_tasks":
		if string(ctx.Method()) == "GET" {
			findTaskHandler(ctx, s)
		}
		return
	// http://localhost:8080/complete_tasks?token=1&topic=php
	case "/complete_tasks":
		if string(ctx.Method()) == "GET" {
			completeTaskHandler(ctx, s)
		}
		return
	// http://localhost:8080/generate_topics?token=1&theme=php
	case "/generate_topics":
		if string(ctx.Method()) == "GET" {
			generateTopicsHandler(ctx, s)
		}
		return
	// http://localhost:8080/create_player?name=john
	case "/create_player":
		if string(ctx.Method()) == "GET" {
			createPlayerHandler(ctx, s)
		}
		return
	default:
		ctx.Error("404 not found.", fasthttp.StatusNotFound)
		return
	}
}

func completeTaskHandler(ctx *fasthttp.RequestCtx, s *taskrpg.Service) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	token := ctx.QueryArgs().Peek(TOKEN)
	topic := ctx.QueryArgs().Peek(TOPIC)
	player, err := s.ValidatePlayerByToken(string(token))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	err = player.CompleteTopic(string(topic))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp := make(map[string]string)
	resp["result"] = fmt.Sprintf("your new level is %v your new xp is %v", player.Level, player.Xp)
	// saving
	s.SetPlayers(player)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	fmt.Fprintf(ctx, string(jsonResp))
	return
}

func createPlayerHandler(ctx *fasthttp.RequestCtx, s *taskrpg.Service) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	name := string(ctx.QueryArgs().Peek(NAME))
	err := s.ValidatePlayerName(name)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp := make(map[string]string)
	pl := s.CreateNewPlayer(name)
	resp["result"] = fmt.Sprintf("Player %s is created with token %s", pl.Name, pl.Token)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	fmt.Fprintf(ctx, string(jsonResp))
	return
}

func generateTopicsHandler(ctx *fasthttp.RequestCtx, s *taskrpg.Service) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	token := ctx.QueryArgs().Peek(TOKEN)
	theme := string(ctx.QueryArgs().Peek(THEME))
	_, err := s.ValidatePlayerByToken(string(token))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	err = s.ValidateTheme(theme)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	question := fmt.Sprintf("Say the most important topics to learn in %s, 10 examples, by 2 words, list separated by commas?", theme)
	topics := normalizeChatGptAnswer(chatgpt.GetChat(question))

	err = s.SaveNewTopics(theme, topics)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}

	resp := make(map[string]string)
	resp["result"] = fmt.Sprintf("finding for theme %s these topics %s", theme, topics)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	fmt.Fprintf(ctx, string(jsonResp))
	return
}

func normalizeChatGptAnswer(s string) string {
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(s, "")
}

func findTaskHandler(ctx *fasthttp.RequestCtx, s *taskrpg.Service) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	token := ctx.QueryArgs().Peek(TOKEN)
	topic := ctx.QueryArgs().Peek(TOPIC)

	pl, err := s.ValidatePlayerByToken(string(token))
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	if len(topic) < 1 {
		ctx.Error("400 bad request", fasthttp.StatusBadRequest)
		return
	}
	resp := make(map[string]string)
	// set a random topic to a player
	curRandTopic, err := s.FindTopic(string(topic))
	if err != nil {
		ctx.Error("400 bad request", fasthttp.StatusBadRequest)
		log.Fatalf("Err: %s", err)
		return
	}
	resp["result"] = fmt.Sprintf("ok for token: %v, topic: %v == %v", token, topic, curRandTopic)
	s.SetTopicAndRemoveOldToPlayer(curRandTopic, pl)

	s.SetPlayers(pl)
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	fmt.Fprintf(ctx, string(jsonResp))
	return
}
