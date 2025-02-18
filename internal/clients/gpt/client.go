package gpt

import (
	"encoding/json"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"rpgMonster/internal/model"
)

const (
	modelName = "gpt-4o"
)

type Client struct {
	restyCl *resty.Client
	token   string
}

func (c *Client) GetCompletion(systemContent, userContent string) (model.GPTAnswer, error) {
	var gptAnswer model.GPTAnswer

	req := c.restyCl.R()
	req.SetHeader("Content-Type", "application/json")
	req.SetHeader("Authorization", "Bearer "+c.token)
	//todo refactoring
	req.SetBody(map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": systemContent,
			},
			{
				"role":    "user",
				"content": userContent,
			},
		},
	})
	resp, err := req.Post("https://api.openai.com/v1/chat/completions")
	if err != nil {
		log.Error().Err(err).Msg("error sending request")
		return gptAnswer, err
	}

	err = json.Unmarshal(resp.Body(), &gptAnswer)
	if err != nil {
		log.Error().Err(err).Msg("error unmarshalling response")
		return gptAnswer, err
	}

	return gptAnswer, nil
}

func New(restyCl *resty.Client) *Client {
	return &Client{
		restyCl: restyCl,
		token:   os.Getenv(model.GPT_TOKEN),
	}
}
