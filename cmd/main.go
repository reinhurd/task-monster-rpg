package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"rpgMonster/internal/clients/dbclient"
	"rpgMonster/internal/clients/gpt"
	"rpgMonster/internal/core"
	"rpgMonster/internal/transport/api"
	"rpgMonster/internal/transport/telegram"

	"github.com/rs/zerolog/log"
)

func main() {
	if err := godotenv.Load("secret.env"); err != nil {
		log.Warn().Msg("secret.env does not exist")
		requiredEnvVars := []string{
			"TG_SECRET_KEY",
			"MONGODB_URI",
		}

		for _, envVar := range requiredEnvVars {
			if os.Getenv(envVar) == "" {
				log.Fatal().Msgf("Required environment variable %s is not set", envVar)
			}
		}
	}

	gptClient := gpt.New(resty.New())
	taskManager := dbclient.NewManager()
	service := core.NewService(gptClient, taskManager)

	tgbot, err := telegram.StartBot(os.Getenv("TG_SECRET_KEY"), true, service)
	if err != nil {
		panic(err)
	}
	router := api.SetupRouter(service)

	//tg bot run in different goroutine
	go func() {
		updChan := tgbot.GetUpdatesChan()
		err = tgbot.HandleUpdate(updChan)
		if err != nil {
			panic(err)
		}
	}()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
