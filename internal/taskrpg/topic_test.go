package taskrpg

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"rpgMonster/models"
	"testing"
)

func Test_topicToCSV(t *testing.T) {
	tests := []struct {
		name   string
		topic  Topic
		expRes []string
	}{
		{
			name: "normal_case",
			topic: Topic{
				MainTheme: "Test",
				Topics:    "Test1,Test2,Test3",
			},
			expRes: []string{"Test", "Test1,Test2,Test3"},
		},
		{
			name: "empty_case",
			topic: Topic{
				MainTheme: "",
				Topics:    "",
			},
			expRes: []string{"", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.topic.ToCSV()
			require.Equal(t, tt.expRes, res)
		})
	}
}

func Test_SaveTopics(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name     string
		topics   []Topic
		mockFunc func(mock *MockIoservice)
	}{
		{
			name: "normal_case",
			topics: []Topic{{
				MainTheme: "Test",
				Topics:    "Test1,Test2,Test3",
			}},
			mockFunc: func(mock *MockIoservice) {
				var expCSV [][]string
				expCSV = append(expCSV, []string{"Test", "Test1,Test2,Test3"})
				mock.EXPECT().SaveTopics(TOPICFILE, expCSV)
			},
		},
		{
			name:     "empty_case",
			topics:   []Topic{},
			mockFunc: func(mock *MockIoservice) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			tt.mockFunc(mock)
			New(mock).SaveTopics(tt.topics)
		})
	}
}

func Test_getTopics(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name     string
		mockFunc func(mock *MockIoservice)
		expRes   []Topic
	}{
		{
			name: "normal_case",
			mockFunc: func(mock *MockIoservice) {
				ret := models.TopicDTO{
					MainTheme: "Test",
					Topics:    "Test1,Test2,Test3",
				}
				mock.EXPECT().GetTopics(TOPICFILE).Return([]models.TopicDTO{ret})
			},
			expRes: []Topic{{
				MainTheme: "Test",
				Topics:    "Test1,Test2,Test3",
			}},
		},
		{
			name: "empty_case",
			mockFunc: func(mock *MockIoservice) {
				mock.EXPECT().GetTopics(TOPICFILE).Return(nil)
			},
			expRes: []Topic{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			tt.mockFunc(mock)
			res := New(mock).getTopics()
			require.Equal(t, tt.expRes, res)
		})
	}
}
