package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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
			res := findRandomTasks(tt.tasks)
			require.Equal(t, tt.expRes, res)
		})
	}
}
