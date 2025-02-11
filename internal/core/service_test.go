package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"rpgMonster/internal/model"
)

func TestNewService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)

	type args struct {
		gptClient GPTClient
		dbClient  DBClient
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "normal_case",
			args: args{
				gptClient: gptClient,
				dbClient:  dbClient,
			},
			want: &Service{
				gptClient: gptClient,
				dbManager: dbClient,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := NewService(tt.args.gptClient, tt.args.dbClient)
			require.Equal(t, tt.want, res)
		})
	}
}

func TestService_CreateTaskFromGPTByRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		req    string
		userID string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient, gpt *MockGPTClient)
		wantErr  bool
		expRes   *model.Task
	}{
		{
			name: "empty_request",
			args: args{
				req:    "",
				userID: "",
			},
			wantErr: true,
		},
		{
			name: "error_case",
			args: args{
				req:    "C#",
				userID: "1",
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient, gpt *MockGPTClient) {
				req := "Write a one single daily task to achieve goal learn C#, in format: 'daily task: task description: requirements to check' and delimiter is comma"
				gpt.EXPECT().GetCompletion(model.GPT_SYSTEM_PROMPT, req).Return(model.GPTAnswer{}, fmt.Errorf("error"))
			},
			expRes: nil,
		},
		{
			name: "normal_case",
			args: args{
				req:    "C#",
				userID: "1",
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient, gpt *MockGPTClient) {
				req := "Write a one single daily task to achieve goal learn C#, in format: 'daily task: task description: requirements to check' and delimiter is comma"
				gpt.EXPECT().GetCompletion(model.GPT_SYSTEM_PROMPT, req).Return(model.GPTAnswer{
					Choices: []model.GPTChoice{
						{
							Message: model.GPTMessage{
								Content: "learn C# and test description",
							},
						},
					},
				}, nil)
				db.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(nil)
			},
			expRes: &model.Task{
				Title:       "learn C#",
				Description: "learn C# and test description",
				Executor:    "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient, gptClient)
			}
			res, err := svc.CreateTaskFromGPTByRequest(tt.args.req, tt.args.userID)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, tt.expRes, res)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, tt.expRes, res)
			}
		})
	}
}

func TestService_CheckPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
		expRes   string
	}{
		{
			name: "normal_case",
			args: args{
				login:    "login",
				password: "password",
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().CheckPassword("login", "password").Return("1", nil)
			},
			expRes: "1",
		},
		{
			name: "error_case",
			args: args{
				login:    "login",
				password: "password",
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().CheckPassword("login", "password").Return("", fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			res, err := svc.CheckPassword(tt.args.login, tt.args.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expRes, res)
			}
		})
	}
}

func TestService_CreateNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
		expRes   string
	}{
		{
			name: "normal_case",
			args: args{
				login:    "login",
				password: "password",
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().CreateNewUser("login", "password").Return("1", nil)
			},
			expRes: "1",
		},
		{
			name: "error_case",
			args: args{
				login:    "login",
				password: "password",
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().CreateNewUser("login", "password").Return("", fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			res, err := svc.CreateNewUser(tt.args.login, tt.args.password)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expRes, res)
			}
		})
	}
}

func TestService_UpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		task *model.Task
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
	}{
		{
			name: "normal_case",
			args: args{
				task: &model.Task{
					BizId:       "1",
					Title:       "title",
					Description: "description",
					Executor:    "1",
				},
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().UpdateTask(gomock.Any(), &model.Task{
					BizId:       "1",
					Title:       "title",
					Description: "description",
					Executor:    "1",
				}).Return(nil)
			},
		},
		{
			name: "error_case",
			args: args{
				task: &model.Task{
					BizId:       "1",
					Title:       "title",
					Description: "description",
					Executor:    "1",
				},
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().UpdateTask(gomock.Any(), &model.Task{
					BizId:       "1",
					Title:       "title",
					Description: "description",
					Executor:    "1",
				}).Return(fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			err := svc.UpdateTask(context.Background(), tt.args.task)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		bizID  string
		userID string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
		expRes   *model.Task
	}{
		{
			name: "normal_case",
			args: args{
				bizID:  "1",
				userID: "1",
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().GetTask(gomock.Any(), "1").Return(&model.Task{
					BizId:    "1",
					Executor: "1",
				}, nil)
			},
			expRes: &model.Task{
				BizId:    "1",
				Executor: "1",
			},
		},
		{
			name: "error_case",
			args: args{
				bizID:  "1",
				userID: "1",
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().GetTask(gomock.Any(), "1").Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			res, err := svc.GetTask(context.Background(), tt.args.bizID, tt.args.userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, tt.expRes, res)
			}
		})
	}
}

func TestService_ValidateUserTG(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		userID int64
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
		expRes   string
	}{
		{
			name: "normal_case",
			args: args{
				userID: 12,
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().GetUserByTGID(int64(12)).Return("1", nil)
			},
			expRes: "1",
		},
		{
			name: "error_case",
			args: args{
				userID: 12,
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().GetUserByTGID(int64(12)).Return("", fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			res, err := svc.ValidateUserTG(tt.args.userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expRes, res)
			}
		})
	}
}

func TestService_ConnectUserToTG(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		userID     string
		telegramID int64
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
	}{
		{
			name: "normal_case",
			args: args{
				userID:     "1",
				telegramID: 12,
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().UpdateUserTGID("1", int64(12)).Return(nil)
			},
		},
		{
			name: "error_case",
			args: args{
				userID:     "1",
				telegramID: 12,
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().UpdateUserTGID("1", int64(12)).Return(fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			err := svc.ConnectUserToTG(tt.args.userID, tt.args.telegramID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_CreateUserFromTG(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		login      string
		password   string
		telegramID int64
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
		expRes   string
	}{
		{
			name: "normal_case",
			args: args{
				login:      "login",
				password:   "password",
				telegramID: 12,
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().CreateNewUserTG("login", "password", int64(12)).Return("1", nil)
			},
			expRes: "1",
		},
		{
			name: "error_case",
			args: args{
				login:      "login",
				password:   "password",
				telegramID: 12,
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().CreateNewUserTG("login", "password", int64(12)).Return("", fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			res, err := svc.CreateUserFromTG(tt.args.login, tt.args.password, tt.args.telegramID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expRes, res)
			}
		})
	}
}

func TestService_CreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		ctx  context.Context
		task *model.Task
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
	}{
		{
			name: "normal_case",
			args: args{
				ctx: context.Background(),
				task: &model.Task{
					BizId: "1",
				},
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().CreateTask(gomock.Any(), &model.Task{
					BizId: "1",
				}).Return(nil)
			},
		},
		{
			name: "error_case",
			args: args{
				ctx: context.Background(),
				task: &model.Task{
					BizId: "1",
				},
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().CreateTask(gomock.Any(), &model.Task{
					BizId: "1",
				}).Return(fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			err := svc.CreateTask(tt.args.ctx, tt.args.task)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestService_GetListTasksByUserID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	type args struct {
		userID string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(db *MockDBClient)
		wantErr  bool
		expRes   []model.Task
	}{
		{
			name: "normal_case",
			args: args{
				userID: "1",
			},
			wantErr: false,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().GetTaskListByUserID("1").Return([]model.Task{
					{
						BizId: "1",
					},
				}, nil)
			},
			expRes: []model.Task{
				{
					BizId: "1",
				},
			},
		},
		{
			name: "error_case",
			args: args{
				userID: "1",
			},
			wantErr: true,
			mockFunc: func(db *MockDBClient) {
				db.EXPECT().GetTaskListByUserID("1").Return(nil, fmt.Errorf("error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockFunc != nil {
				tt.mockFunc(dbClient)
			}
			res, err := svc.GetListTasksByUserID(context.Background(), tt.args.userID)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expRes, res)
			}
		})
	}
}

func TestService_GetTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	dbClient := NewMockDBClient(ctrl)
	gptClient := NewMockGPTClient(ctrl)
	svc := NewService(gptClient, dbClient)

	tests := []struct {
		name string
	}{
		{
			name: "normal_case",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := svc.GetTemplate()
			require.NotEmpty(t, res)
		})
	}
}
