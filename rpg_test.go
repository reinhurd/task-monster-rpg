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

func Test_completeTopic(t *testing.T) {
	tests := []struct {
		name   string
		player Player
		topic  string
		expRes int64
		expErr func(err error)
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
			topic:  "PHP",
			expRes: 110,
			expErr: func(err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "invalid_topic",
			player: Player{
				Name:        "",
				Token:       "",
				CurrentTask: "",
				Level:       0,
				Xp:          0,
				Health:      0,
			},
			topic:  "non",
			expRes: 0,
			expErr: func(err error) {
				require.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.player.completeTopic(tt.topic)
			tt.expErr(err)
			require.Equal(t, tt.expRes, tt.player.Xp)
		})
	}
}

func Test_setTopicAndRemoveOldToPlayer(t *testing.T) {
	tests := []struct {
		name     string
		player   Player
		newtopic string
		expRes   int64
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
			newtopic: "PHP2",
			expRes:   80,
		},
		{
			name: "invalid_topic",
			player: Player{
				Name:        "",
				Token:       "",
				CurrentTask: "",
				Level:       0,
				Xp:          0,
				Health:      0,
			},
			newtopic: "",
			expRes:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTopicAndRemoveOldToPlayer(tt.newtopic, &tt.player)
			require.Equal(t, tt.expRes, tt.player.Xp)
		})
	}
}

func Test_generateToken(t *testing.T) {
	tests := []struct {
		name   string
		expRes func(res string)
	}{
		{
			name: "normal_case",
			expRes: func(res string) {
				require.NotEmpty(t, res)
				require.Equal(t, DEFAULT_TOKEN_LENGHT, len(res))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := generateToken()
			tt.expRes(res)
		})
	}
}

func Test_stringToInt(t *testing.T) {
	tests := []struct {
		name   string
		req    string
		expRes int64
	}{
		{
			name:   "normal_case",
			req:    "1234",
			expRes: 1234,
		},
		{
			name:   "normal_case_2",
			req:    "01234",
			expRes: 1234,
		},
		{
			name:   "invalid_case",
			req:    "abd4332",
			expRes: 0,
		},
		{
			name:   "invalid_case_2",
			req:    "abc",
			expRes: 0,
		},
		{
			name:   "empty_case",
			req:    "",
			expRes: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := stringToInt(tt.req)
			require.Equal(t, tt.expRes, res)
		})
	}
}
