package taskrpg

import (
	"github.com/stretchr/testify/require"
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
