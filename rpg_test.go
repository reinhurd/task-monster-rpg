package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_toCSV(t *testing.T) {
	tests := []struct {
		name   string
		player Player
		expRes []string
	}{
		{
			name: "normal_case",
			player: Player{
				Name:        "PersonOne",
				Token:       "123456",
				CurrentTask: "PHP",
				Level:       1,
				Xp:          100,
				Health:      100,
			},
			expRes: []string{"PersonOne", "123456", "PHP", "1", "100", "100"},
		},
		{
			name: "empty_case",
			player: Player{
				Name:        "",
				Token:       "",
				CurrentTask: "",
				Level:       0,
				Xp:          0,
				Health:      0,
			},
			expRes: []string{"", "", "", "0", "0", "0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.player.toCSV()
			require.Equal(t, tt.expRes, res)
		})
	}
}

func Test_setNewLevel(t *testing.T) {
	tests := []struct {
		name     string
		player   Player
		expResXp int64
		expResLv int64
	}{
		{
			name: "normal_case",
			player: Player{
				Name:        "PersonOne",
				Token:       "123456",
				CurrentTask: "PHP",
				Level:       1,
				Xp:          100,
				Health:      100,
			},
			expResXp: 100,
			expResLv: 1,
		},
		{
			name: "normal_case_2",
			player: Player{
				Name:        "PersonOne",
				Token:       "123456",
				CurrentTask: "PHP",
				Level:       1,
				Xp:          1000,
				Health:      100,
			},
			expResXp: 0,
			expResLv: 2,
		},
		{
			name: "empty_case",
			player: Player{
				Name:        "",
				Token:       "",
				CurrentTask: "",
				Level:       0,
				Xp:          0,
				Health:      0,
			},
			expResXp: 0,
			expResLv: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.setNewLevel()
			require.Equal(t, tt.expResXp, tt.player.Xp)
			require.Equal(t, tt.expResLv, tt.player.Level)
		})
	}
}

func Test_completeTasksForXp(t *testing.T) {
	tests := []struct {
		name   string
		player Player
		expRes int64
	}{
		{
			name: "normal_case",
			player: Player{
				Name:        "PersonOne",
				Token:       "123456",
				CurrentTask: "PHP",
				Level:       1,
				Xp:          100,
				Health:      100,
			},
			expRes: 110,
		},
		{
			name: "null_case",
			player: Player{
				Name:        "",
				Token:       "",
				CurrentTask: "",
				Level:       0,
				Xp:          0,
				Health:      0,
			},
			expRes: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.player.completeTasksForXp(DEFAULT_REWARD)
			require.Equal(t, tt.expRes, tt.player.Xp)
		})
	}
}
