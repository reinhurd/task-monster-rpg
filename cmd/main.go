package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"rpgMonster/internal/clients/gpt"
	"rpgMonster/internal/clients/telegram"
	"rpgMonster/internal/core"
	"rpgMonster/internal/model"
	"rpgMonster/internal/transport"

	"github.com/joho/godotenv"

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

	router := transport.SetupRouter()

	tgbot, err := telegram.StartBot(os.Getenv("TG_SECRET_KEY"), true)

	if err != nil {
		panic(err)
	}

	// gptClient := gpt.New(resty.New())
	// p1 := core.GeneratePlayer()
	// exampleWayOfUser(gptClient, &p1)

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
	if err := srv.Shutdown(ctx); err != nil {
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

func exampleWayOfUser(gpt *gpt.Client, p1 *model.Player) {
	//todo clear "sure" and other not csv stuff
	//todo separate daily tasks
	//if player.json doesn't exist
	if _, err := os.Stat("player.json"); os.IsNotExist(err) {
		resp, err := gpt.GetCompletion("You a personal assistant, helping people to set concrete detailed steps to achieve goals", "I want to learn php")
		if err != nil {
			log.Error().Err(err).Msg("error getting completion")
			return
		}

		//set goal
		p1.Goal = "learn PHP"
		p1.GoalDetails = append(p1.GoalDetails, resp.Choices[0].Message.Content)
		resp, err = gpt.GetCompletion("You a personal assistant, helping people to set concrete detailed steps to achieve goals", "Write a daily tasks to achieve goal learn PHP"+
			", in format: 'daily task: task description' and delimiter is comma")
		if err != nil {
			log.Error().Err(err).Msg("error getting completion")
			return
		}
		var dailyMap = make(map[string]model.Daily)
		dailyMap["php"] = model.Daily{Task: resp.Choices[0].Message.Content, Completed: false}
		p1.Dailies = dailyMap
		os.WriteFile("player.json", []byte(fmt.Sprintf("%+v", p1)), 0644)
	}
	playerFile, err := os.ReadFile("player.json")
	if err != nil {
		log.Error().Err(err).Msg("error reading player.json")
		return
	}
	//read json from file
	_, err = fmt.Sscanf(string(playerFile), "%+v", p1)

	//battle engine run regularly - when the player set goal
	go func() {
		//battleEngine
		for {
			victCount := 0
			m := core.GenerateMonster(1)
			w := core.Battle(p1, m)
			if w {
				victCount++
				log.Info().Msgf("Player won! XP: %d, Victories: %d", p1.CurrentXP, victCount)
			} else {
				log.Info().Msgf("Monster won!")
				p1.CurrentXP = 0
			}
		}
	}()

	//if player call api to complete daily, he became stronger

	//if player get another level, he gets another daily and another info on his goal

	//show player progress and his dailies on tg and web
}
