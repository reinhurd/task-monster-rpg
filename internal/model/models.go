package model

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

type Player struct {
	Name        string
	HP          int
	Atk         int
	CurrentXP   int
	Level       int
	Goal        string
	GoalDetails []string         //to add some details as goal progressed - this is generally reward
	Dailies     map[string]Daily //as daily progressed, you became stronger
}

type Daily struct {
	Task      string
	Completed bool
}

type Monster struct {
	Name string
	HP   int
	Atk  int
	XP   int
}
