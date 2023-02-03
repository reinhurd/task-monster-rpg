package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"
)

var headerContentTypeJson = []byte("application/json")

var client *fasthttp.Client

const (
	CHATGPT_API_URL     = "https://api.openai.com/v1/completions"
	CHATGRPMODEL        = "text-davinci-003"
	DEFAULT_MAX_TOKENS  = 4000
	DEFAULT_TEMPERATURE = 1.0
)

type ChatGPTEntity struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type ChatGPTResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
	Usage   Usage     `json:"usage"`
}

type Choices struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func getChat(question string) string {
	readTimeout, _ := time.ParseDuration("10000ms")
	writeTimeout, _ := time.ParseDuration("10000ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	client = &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	if question == "" {
		question = "What is php?"
	}

	reqEntity := &ChatGPTEntity{
		Model:       CHATGRPMODEL,
		Prompt:      question,
		MaxTokens:   DEFAULT_MAX_TOKENS,
		Temperature: DEFAULT_TEMPERATURE,
	}
	reqEntityBytes, _ := json.Marshal(reqEntity)

	reqTimeout := time.Duration(10000) * time.Millisecond

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(CHATGPT_API_URL)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes(headerContentTypeJson)
	req.Header.Set("Authorization", PRIVATE_TOKEN)
	req.SetBodyRaw(reqEntityBytes)
	resp := fasthttp.AcquireResponse()
	err := client.DoTimeout(req, resp, reqTimeout)
	fasthttp.ReleaseRequest(req)
	if err == nil {
		statusCode := resp.StatusCode()
		respBody := resp.Body()
		//fmt.Printf("DEBUG Response: %s\n", respBody)
		if statusCode == http.StatusOK {
			respEntity := &ChatGPTResponse{}
			err = json.Unmarshal(respBody, respEntity)
			if err == io.EOF || err == nil {
				//fmt.Printf("DEBUG Parsed Response: %v\n", respEntity)
				fmt.Println(respEntity.Choices[0].Text)
				return respEntity.Choices[0].Text
			} else {
				fmt.Fprintf(os.Stderr, "ERR failed to parse reponse: %v\n", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "ERR invalid HTTP response code: %d\n", statusCode)
		}
	} else {
		errName, known := httpConnError(err)
		if known {
			fmt.Fprintf(os.Stderr, "WARN conn error: %v\n", errName)
		} else {
			fmt.Fprintf(os.Stderr, "ERR conn failure: %v %v\n", errName, err)
		}
	}
	fasthttp.ReleaseResponse(resp)
	return ""
}

func httpConnError(err error) (string, bool) {
	errName := ""
	known := false
	if err == fasthttp.ErrTimeout {
		errName = "timeout"
		known = true
	} else if err == fasthttp.ErrNoFreeConns {
		errName = "conn_limit"
		known = true
	} else if err == fasthttp.ErrConnectionClosed {
		errName = "conn_close"
		known = true
	} else {
		errName = reflect.TypeOf(err).String()
		if errName == "*net.OpError" {
			// Write and Read errors are not so often and in fact they just mean timeout problems
			errName = "timeout"
			known = true
		}
	}
	return errName, known
}
