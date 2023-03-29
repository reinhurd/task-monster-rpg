package taskrpg

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	ios "rpgMonster/internal/ioservice"
	"rpgMonster/models"
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

func Test_toPlayer(t *testing.T) {
	tests := []struct {
		name       string
		playersDTO []models.PlayerDTO
		players    []Player
	}{
		{
			name: "normal_case",
			playersDTO: []models.PlayerDTO{{
				Name:        "PersonOne",
				Token:       "123456",
				CurrentTask: "PHP",
				Level:       "1",
				Xp:          "100",
				Health:      "100",
			},
			},
			players: []Player{{
				Name:        "PersonOne",
				Token:       "123456",
				CurrentTask: "PHP",
				Level:       1,
				Xp:          100,
				Health:      100,
			},
			},
		},
		{
			name: "empty_case",
			playersDTO: []models.PlayerDTO{{
				Name:        "",
				Token:       "",
				CurrentTask: "",
				Level:       "",
				Xp:          "",
				Health:      "",
			},
			},
			players: []Player{{
				Name:        "",
				Token:       "",
				CurrentTask: "",
				Level:       0,
				Xp:          0,
				Health:      0,
			},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := toPlayer(tt.playersDTO)
			require.Equal(t, tt.players, res)
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
			err := tt.player.CompleteTopic(tt.topic)
			tt.expErr(err)
			require.Equal(t, tt.expRes, tt.player.Xp)
		})
	}
}

func Test_setTopicAndRemoveOldToPlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
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
			s := New(NewMockIoservice(ctrl))
			s.SetTopicAndRemoveOldToPlayer(tt.newtopic, &tt.player)
			require.Equal(t, tt.expRes, tt.player.Xp)
		})
	}
}

func Test_generateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
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
			s := New(NewMockIoservice(ctrl))
			res := s.generateToken()
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

func Test_findRandomTasks(t *testing.T) {
	tests := []struct {
		name   string
		tasks  string
		expRes string
	}{
		{
			name:   "empty_string",
			tasks:  "",
			expRes: "",
		},
		{
			name:   "all_correct",
			tasks:  "test,test,test",
			expRes: "test",
		},
		{
			name:   "one_value",
			tasks:  "one_value",
			expRes: "one_value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iosMock := ios.New()
			s := New(iosMock)
			res := s.findRandomTasks(tt.tasks)
			require.Equal(t, tt.expRes, res)
		})
	}
}

func Test_loadPlayers(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name     string
		mockFunc func(mock *MockIoservice)
		expRes   []Player
	}{
		{
			name: "normal_case",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().
					LoadPlayers(PLAYERFILE).
					Return([]models.PlayerDTO{{Name: "1", Token: "1", CurrentTask: "1", Level: "1", Xp: "1", Health: "1"}})
			},
			expRes: []Player{{
				Name:        "1",
				Token:       "1",
				CurrentTask: "1",
				Level:       1,
				Xp:          1,
				Health:      1,
			}},
		},
		{
			name: "empty_case",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().
					LoadPlayers(PLAYERFILE).
					Return(nil)
			},
			expRes: []Player{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			s := New(mock)
			tt.mockFunc(mock)
			res := s.loadPlayers()
			require.Equal(t, tt.expRes, res)
		})
	}
}

func Test_savePlayers(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name     string
		req      []Player
		mockFunc func(mock *MockIoservice)
	}{
		{
			name: "normal_case",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().
					SavePlayers(PLAYERFILE, [][]string{{"name", "token", "task", "level", "xp", "health"}, {"1", "1", "1", "1", "1", "1"}}).
					Return()
			},
			req: []Player{{
				Name:        "1",
				Token:       "1",
				CurrentTask: "1",
				Level:       1,
				Xp:          1,
				Health:      1,
			}},
		},
		{
			name: "empty_case",
			req:  []Player{},
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().
					SavePlayers(PLAYERFILE, [][]string{{"name", "token", "task", "level", "xp", "health"}}).
					Return()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			s := New(mock)
			tt.mockFunc(mock)
			s.SavePlayers(tt.req)
		})
	}
}

func Test_ValidatePlayerName(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name     string
		req      string
		mockFunc func(mock *MockIoservice)
		expErr   error
	}{
		{
			name: "player_exists",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().LoadPlayers(PLAYERFILE).Return([]models.PlayerDTO{{Name: "TeST", Token: "1", CurrentTask: "1", Level: "1", Xp: "1", Health: "1"}})
			},
			req:    "TEst",
			expErr: fmt.Errorf("player name %s already exists", "TEst"),
		},
		{
			name: "normal_case",
			req:  "test",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().LoadPlayers(PLAYERFILE).Return(nil)
			},
			expErr: nil,
		},
		{
			name:     "empty_name",
			req:      "",
			mockFunc: func(mock *MockIoservice) {},
			expErr:   errors.New("player name is empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			s := New(mock)
			tt.mockFunc(mock)
			err := s.ValidatePlayerName(tt.req)
			require.Equal(t, tt.expErr, err)
		})
	}
}

func Test_CreateNewPlayer(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name     string
		req      string
		mockFunc func(mock *MockIoservice)
		expRes   Player
	}{
		{
			name: "player_exists",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().LoadPlayers(PLAYERFILE).Return([]models.PlayerDTO{{Name: "TeST", Token: "1", CurrentTask: "1", Level: "1", Xp: "1", Health: "1"}})
				mock.EXPECT().SavePlayers(PLAYERFILE, gomock.Any()).Return()
			},
			req:    "TEst",
			expRes: Player{Name: "TEst", CurrentTask: "", Level: 1, Xp: 0, Health: 0},
		},
		{
			name: "empty_case",
			req:  "test",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().LoadPlayers(PLAYERFILE).Return(nil)
				mock.EXPECT().SavePlayers(PLAYERFILE, gomock.Any()).Return()
			},
			expRes: Player{Name: "test", CurrentTask: "", Level: 1, Xp: 0, Health: 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			s := New(mock)
			tt.mockFunc(mock)
			res := s.CreateNewPlayer(tt.req)
			require.Equal(t, tt.expRes.Name, res.Name)
			require.Equal(t, tt.expRes.CurrentTask, res.CurrentTask)
			require.Equal(t, tt.expRes.Level, res.Level)
			require.Equal(t, tt.expRes.Xp, res.Xp)
			require.Equal(t, tt.expRes.Health, res.Health)
		})
	}
}

func Test_ValidatePlayerByToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	tests := []struct {
		name     string
		req      string
		mockFunc func(mock *MockIoservice)
		expRes   *Player
		expErr   error
	}{
		{
			name: "player_exists",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().LoadPlayers(PLAYERFILE).Return([]models.PlayerDTO{{Name: "TeST", Token: "12345", CurrentTask: "1", Level: "1", Xp: "1", Health: "1"}})
			},
			req:    "12345",
			expRes: &Player{Name: "TeST", Token: "12345", CurrentTask: "1", Level: 1, Xp: 1, Health: 1},
		},
		{
			name: "empty_case",
			req:  "43214",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().LoadPlayers(PLAYERFILE).Return([]models.PlayerDTO{{Name: "TeST", Token: "12345", CurrentTask: "1", Level: "1", Xp: "1", Health: "1"}})
			},
			expRes: nil,
			expErr: errors.New("no token found for players"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			s := New(mock)
			tt.mockFunc(mock)
			player, err := s.ValidatePlayerByToken(tt.req)
			require.Equal(t, tt.expErr, err)
			if err == nil {
				require.Equal(t, tt.expRes.CurrentTask, player.CurrentTask)
				require.Equal(t, tt.expRes.Level, player.Level)
				require.Equal(t, tt.expRes.Xp, player.Xp)
				require.Equal(t, tt.expRes.Health, player.Health)
				require.Equal(t, tt.expRes.Token, player.Token)
			}
		})
	}
}

func Test_ValidateTheme(t *testing.T) {
	tests := []struct {
		name   string
		theme  string
		expErr error
	}{
		{
			name:   "empty_string",
			theme:  "",
			expErr: errors.New("invalid theme"),
		},
		{
			name:  "all_correct",
			theme: "test,test,test",
		},
		{
			name:   "single_value",
			theme:  "1",
			expErr: errors.New("invalid theme"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iosMock := ios.New()
			s := New(iosMock)
			err := s.ValidateTheme(tt.theme)
			require.Equal(t, tt.expErr, err)
		})
	}
}
