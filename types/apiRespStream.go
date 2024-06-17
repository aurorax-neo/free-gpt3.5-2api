package types

import (
	"encoding/json"
	"strings"
	"time"
)

// ApiRespStream  represents the JSON structure
type ApiRespStream struct {
	ID      string            `json:"id,omitempty"`
	Object  string            `json:"object,omitempty"`
	Created int64             `json:"created,omitempty"`
	Model   string            `json:"model,omitempty"`
	Choices []ApiStreamChoice `json:"choices,omitempty"`
}

// ApiStreamChoice represents the nested "choices" object in the JSON
type ApiStreamChoice struct {
	Delta        ApiStreamDelta `json:"delta,omitempty"`
	Index        int            `json:"index,omitempty"`
	FinishReason interface{}    `json:"finish_reason,omitempty"`
}

// ApiStreamDelta represents the nested "delta" object in the JSON
type ApiStreamDelta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

func NewApiRespStream(id string, model string, content string) *ApiRespStream {
	// 生成响应 model
	apiRespStream := &ApiRespStream{
		ID:      id,
		Created: time.Now().Unix(),
		Object:  "chat.completion.chunk",
		Model:   model,
		Choices: []ApiStreamChoice{
			{
				Delta: ApiStreamDelta{
					Content: content,
				},
				Index:        0,
				FinishReason: nil,
			},
		},
	}
	return apiRespStream
}

func ConvertToString(id string, model string, chatResp *ChatResp, previousText *StringStruct, role bool) string {
	apiRespJson := NewApiRespStream(id, model, strings.Replace(chatResp.Message.Content.Parts[0].(string), previousText.Text, "", 1))
	if role {
		apiRespJson.Choices[0].Delta.Role = chatResp.Message.Author.Role
	} else if apiRespJson.Choices[0].Delta.Content == "" || (strings.HasPrefix(chatResp.Message.Metadata.ModelSlug, "gpt-4") && apiRespJson.Choices[0].Delta.Content == "【") {
		return apiRespJson.Choices[0].Delta.Content
	}
	previousText.Text = chatResp.Message.Content.Parts[0].(string)
	data, _ := json.Marshal(apiRespJson)
	return "data: " + string(data) + "\n\n"
}

func StopChunk(id string, model string, finishReason string) ApiRespStream {
	return ApiRespStream{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []ApiStreamChoice{
			{
				Index:        0,
				FinishReason: finishReason,
			},
		},
	}
}

func (ARS *ApiRespStream) String() string {
	resp, _ := json.Marshal(ARS)
	return string(resp)
}
