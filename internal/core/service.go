package core

import (
	"math/rand"
	"time"

	"github.com/rs/zerolog/log"
)

type Player struct {
	Name string
	HP   int
	Atk  int
}

func GeneratePlayers() (p1, p2 Player) {
	p1 = Player{"Player1", rand.Intn(100), rand.Intn(10)}
	p2 = Player{"Player2", rand.Intn(100), rand.Intn(10)}

	return
}

func Battle(p1, p2 Player) string {
	//each player attacks the other until one of them dies in 2 seconds
	for p1.HP > 0 && p2.HP > 0 {
		p1.HP -= rand.Intn(p2.Atk)
		p2.HP -= rand.Intn(p1.Atk)

		log.Info().Msgf("%s HP: %d, %s HP: %d", p1.Name, p1.HP, p2.Name, p2.HP)

		wait := time.After(2 * time.Second)
		<-wait
	}

	if p1.HP > p2.HP {
		return p1.Name
	}

	return p2.Name
}
