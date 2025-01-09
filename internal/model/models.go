package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GPTAnswer struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
}

type User struct {
	BizID     string `bson:"biz_id"`
	Login     string `bson:"login"`
	Password  string `bson:"password"`
	Salt      string `bson:"salt"`
	StrikeDay int    `bson:"strike_day"`
}

type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	BizId       string             `bson:"biz_id"`
	Title       string             `bson:"title"`
	Description string             `bson:"description"`
	Executor    string             `bson:"executor"` // ID of the user executing the task
	Reviewer    *string            `bson:"reviewer"` // Optional ID of the reviewing user
	Completed   bool               `bson:"completed"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	Deadline    primitive.DateTime `bson:"deadline"`
	Tags        []string           `bson:"tags"`
}
