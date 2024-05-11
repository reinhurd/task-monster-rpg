package core

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"

	"rpgMonster/internal/model"
)

func randomInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + rand.Intn(max)
}

func GeneratePlayer() model.Player {
	//generate random int from 1 to 10
	p := model.Player{Name: "Player", HP: randomInt(1, 10), Atk: randomInt(1, 10), Level: 1}

	return p
}

func GenerateMonster(difficulty int) (m model.Monster) {
	if difficulty < 1 {
		difficulty = 1
	}
	m = model.Monster{Name: "Monster", HP: randomInt(1, 10) * difficulty, Atk: randomInt(1, 10) * difficulty, XP: randomInt(1, 10) * difficulty}

	return
}

func AddXP(p *model.Player, xp int) {
	p.CurrentXP += xp
	if p.CurrentXP >= p.Level*10 {
		p.Level++
		p.HP += 10
		p.Atk += 5
	}
	log.Info().Msgf("LevelUp! %s, HP: %d, Atk: %d, XP: %d, Level: %d", p.Name, p.HP, p.Atk, p.CurrentXP, p.Level)
}

//todo add some weapons and armor?

func Battle(p1 *model.Player, m model.Monster) bool {
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
		AddXP(p1, m.XP)
		return true
	}

	return false
}
