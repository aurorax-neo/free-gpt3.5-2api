package respModel

import "time"

type ApiRespJson struct {
	ID      string              `json:"id"`
	Object  string              `json:"object"`
	Created int64               `json:"created"`
	Model   string              `json:"model"`
	Usage   ApiRespJsonUsage    `json:"usage"`
	Choices []ApiRespJsonChoice `json:"choices"`
}

type ApiRespJsonMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiRespJsonChoice struct {
	Message      ApiRespJsonMessage `json:"message"`
	FinishReason string             `json:"finish_reason"`
	Index        int                `json:"index"`
}

type ApiRespJsonUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewApiRespJson(id string, model string, content string) *ApiRespJson {
	apiRespObj := &ApiRespJson{
		ID:      id,
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