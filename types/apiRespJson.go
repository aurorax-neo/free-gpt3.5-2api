package types

import (
	"time"
)

type ApiRespJson struct {
	ID      string              `json:"id,omitempty"`
	Object  string              `json:"object,omitempty"`
	Created int64               `json:"created,omitempty"`
	Model   string              `json:"model,omitempty"`
	Usage   ApiRespJsonUsage    `json:"usage,omitempty"`
	Choices []ApiRespJsonChoice `json:"choices,omitempty"`
}

type ApiRespJsonMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ApiRespJsonChoice struct {
	Delta        ApiRespJsonChoiceDelta `json:"delta,omitempty"`
	Message      ApiRespJsonMessage     `json:"message,omitempty"`
	FinishReason string                 `json:"finish_reason,omitempty"`
	Index        int                    `json:"index,omitempty"`
}

type ApiRespJsonChoiceDelta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

type ApiRespJsonUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewApiRespJson(model string, content string) *ApiRespJson {
	apiRespObj := &ApiRespJson{
		ID:      GenerateID(29),
		Created: time.Now().Unix(),
		Object:  "chat.completion",
		Model:   model,
		Usage: ApiRespJsonUsage{
			PromptTokens:     0,
			CompletionTokens: 0,
			TotalTokens:      0,
		},
		Choices: []ApiRespJsonChoice{
			{
				Message: ApiRespJsonMessage{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
				Index:        0,
			},
		},
	}
	return apiRespObj
}
