package main

import (
	"fmt"
	"log"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"

	"rpgMonster/internal/http_handler"
	"rpgMonster/internal/ioservice"
	"rpgMonster/internal/taskrpg"
	"rpgMonster/internal/tgbot"
)

func main() {
	// todo add context
	initApp()
	fmt.Printf("Starting server...\n")

	ln, err := reuseport.Listen("tcp4", "localhost:8080")
	if err != nil {
		log.Fatalf("error in reuseport listener: %v", err)
	}

	if err = fasthttp.Serve(ln, http_handler.Handler); err != nil {
		log.Fatalf("error in fasthttp Server: %v", err)
	}
}

func initApp() {
	ios := ioservice.New()
	s := taskrpg.New(ios)
	s.SaveTopics(taskrpg.DEFAULT_TOPICS)
	s.SavePlayers(taskrpg.DEFAULT_PLAYERS_DATA)

	tgbot.StartBot()

	return
}
