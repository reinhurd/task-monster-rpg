package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"rpgMonster/internal/clients/telegram"
	"rpgMonster/internal/core"
	"rpgMonster/internal/transport"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := godotenv.Load("secret.env"); err != nil {
		panic(err)
	}

	r := transport.SetupRouter()
	tgbot, err := telegram.StartBot(os.Getenv("TG_SECRET_KEY"), true)
	if err != nil {
		panic(err)
	}

	//gptClient := gpt.New(resty.New())
	//resp, err := gptClient.GetCompletion()
	//fmt.Println(resp.Choices[0].Message.Content)

	go func() {
		updChan := tgbot.GetUpdatesChan()
		err = tgbot.HandleUpdate(updChan)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		//battleEngine
		for {
			victCount := 0
			p1 := core.GeneratePlayer()
			m := core.GenerateMonster(1)
			w := core.Battle(&p1, m)
			if w {
				victCount++
				log.Info().Msgf("Player won! XP: %d, Victories: %d", p1.CurrentXP, victCount)
				for i := 1; i < 5; i++ {
					m = core.GenerateMonster(i)
					w = core.Battle(&p1, m)
					if w {
						victCount++
						log.Info().Msgf("Player won! XP: %d, Victories: %d", p1.CurrentXP, victCount)
					} else {
						log.Info().Msgf("Monster won!")
						break
					}
				}
			} else {
				log.Info().Msgf("Monster won!")
			}
		}
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err)
		}
	}()

	quit := make(chan os.Signal, 2)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = srv.Shutdown(ctx); err != nil {
		_, errTg := tgbot.SendToLastChat("Service is shutting down with error")
		if errTg != nil {
			log.Info().Msgf("Error sending to telegram: %v", errTg)
		}
		log.Fatal().Err(err)
	}

	<-ctx.Done()
	log.Info().Msg("timeout of 5 seconds")
	_, errTg := tgbot.SendToLastChat("Service is shutting down by timeout")
	if errTg != nil {
		log.Info().Msgf("Error sending to telegram: %v", errTg)
	}
	log.Info().Msg("Server exiting")
}
