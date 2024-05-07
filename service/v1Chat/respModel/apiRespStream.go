package respModel

import "time"

// ApiRespStream  represents the JSON structure
type ApiRespStream struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []ApiStreamChoice `json:"choices"`
}

// ApiStreamChoice represents the nested "choices" object in the JSON
type ApiStreamChoice struct {
	Delta        ApiStreamDelta `json:"delta"`
	Index        int            `json:"index"`
	FinishReason string         `json:"finish_reason"`
}

// ApiStreamDelta represents the nested "delta" object in the JSON
type ApiStreamDelta struct {
	Content string `json:"content"`
}

func NewApiRespStream(id string, model string, content string, finishReason string) *ApiRespStream {
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
				FinishReason: finishReason,
			},
		},
	}
	return apiRespStream
}
