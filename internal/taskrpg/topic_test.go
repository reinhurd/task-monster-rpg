package taskrpg

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"rpgMonster/models"
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

func Test_SaveNewTopics(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name          string
		theme, topics string
		mockFunc      func(mock *MockIoservice)
		expErr        error
	}{
		{
			name:   "normal_case",
			theme:  "TestNew",
			topics: "TestNew1,TestNew2,TestNew3",
			mockFunc: func(mock *MockIoservice) {
				ret := models.TopicDTO{
					MainTheme: "Test",
					Topics:    "Test1,Test2,Test3",
				}
				mock.EXPECT().GetTopics(TOPICFILE).Return([]models.TopicDTO{ret})
				mock.EXPECT().SaveTopics(TOPICFILE, [][]string{{"Test", "Test1,Test2,Test3"}, {"testnew", "testnew1,testnew2,testnew3"}})
			},
		},
		{
			name:   "replace_case",
			theme:  "Test",
			topics: "TestNew1,TestNew2,TestNew3",
			mockFunc: func(mock *MockIoservice) {
				ret := models.TopicDTO{
					MainTheme: "Test",
					Topics:    "Test1,Test2,Test3",
				}
				mock.EXPECT().GetTopics(TOPICFILE).Return([]models.TopicDTO{ret})
				mock.EXPECT().SaveTopics(TOPICFILE, [][]string{{"test", "testnew1,testnew2,testnew3"}})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			tt.mockFunc(mock)
			err := New(mock).SaveNewTopics(tt.theme, tt.topics)
			require.Equal(t, tt.expErr, err)
		})
	}
}

func Test_makeTopicsAsMap(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name   string
		req    []Topic
		expRes map[string]string
	}{
		{
			name:   "normal_case",
			req:    DEFAULT_TOPICS,
			expRes: map[string]string{"golang": "Concurrency,Parallelism,Goroutine,Frameworks", "php": "Concurrency,Parallelism,PHP9,Frameworks"},
		},
		{
			name:   "empty_case",
			req:    []Topic{},
			expRes: map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := NewMockIoservice(ctrl)
			res := New(mock).makeTopicsAsMap(tt.req)
			require.Equal(t, tt.expRes, res)
		})
	}
}
