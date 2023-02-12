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
