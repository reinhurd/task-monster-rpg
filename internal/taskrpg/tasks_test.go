package taskrpg

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"rpgMonster/models"
)

func Test_FindTopic(t *testing.T) {
	ctrl := gomock.NewController(t)
	iosMock := NewMockIoservice(ctrl)

	tests := []struct {
		name       string
		topic      string
		mockFunc   func(mock *MockIoservice)
		expErrFunc func(res string, err error)
	}{
		{
			name:  "success",
			topic: "Test",
			mockFunc: func(mock *MockIoservice) {
				ret := models.TopicDTO{
					MainTheme: "Test",
					Topics:    "Test1,Test2,Test3",
				}
				mock.EXPECT().GetTopics(TOPICFILE).Return([]models.TopicDTO{ret})
			},
			expErrFunc: func(res string, err error) {
				require.NoError(t, err)
				isRandTopic := res == "Test1" || res == "Test2" || res == "Test3"
				require.True(t, isRandTopic)
			},
		},
		{
			name:  "not_found",
			topic: "Test66",
			mockFunc: func(mock *MockIoservice) {
				ret := models.TopicDTO{
					MainTheme: "Test",
					Topics:    "Test1,Test2,Test3",
				}
				mock.EXPECT().GetTopics(TOPICFILE).Return([]models.TopicDTO{ret})
			},
			expErrFunc: func(res string, err error) {
				require.Error(t, err)
				require.Equal(t, "Test66 not found", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(iosMock)
			tt.mockFunc(iosMock)
			res, err := s.FindTopic(tt.topic)
			tt.expErrFunc(res, err)
		})
	}
}
