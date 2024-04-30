package core

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

type Player struct {
	Name      string
	HP        int
	Atk       int
	CurrentXP int
	Level     int
}

//todo add level up logic - which level is x * xp

type Monster struct {
	Name string
	HP   int
	Atk  int
	XP   int
}

func randomInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.Intn(max)
}

func GeneratePlayer() Player {
	//generate random int from 1 to 10
	p := Player{"Player", randomInt(1, 10), randomInt(1, 10), 0, 1}

	return p
}

func GenerateMonster(difficulty int) (m Monster) {
	m = Monster{"Monster", randomInt(1, 10) * difficulty, randomInt(1, 10) * difficulty, randomInt(1, 10) * difficulty}

	return
}

//todo add some weapons and armor?

func Battle(p1 *Player, m Monster) bool {
	log.Info().Msgf("Player: %s, HP: %d, Atk: %d, XP: %d, Level: %d", p1.Name, p1.HP, p1.Atk, p1.CurrentXP, p1.Level)
	log.Info().Msgf("Monster: %s, HP: %d, Atk: %d, XP: %d", m.Name, m.HP, m.Atk, m.XP)
	//each player attacks the other until one of them dies in 2 seconds
	for p1.HP > 0 && m.HP > 0 {
		//todo add some initiative logic
		p1.HP -= randomInt(1, m.Atk)
		m.HP -= randomInt(1, p1.Atk)

		//log.Info().Msgf("%s HP: %d, %s HP: %d", p1.Name, p1.HP, m.Name, m.HP)

		wait := time.After(2 * time.Second)
		<-wait
	}

	if p1.HP > m.HP {
		//did the player win?
		log.Info().Msgf("%s won!", p1.Name)
		p1.CurrentXP += m.XP
		return true
	}

	return false
}
