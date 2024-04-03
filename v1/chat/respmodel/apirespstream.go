package respmodel

// StreamObj ApiRespStream represents the JSON structure
type StreamObj struct {
	ID      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []StreamChoiceObj `json:"choices"`
}

// StreamChoiceObj represents the nested "choices" object in the JSON
type StreamChoiceObj struct {
	Delta        StreamDeltaObj `json:"delta"`
	Index        int            `json:"index"`
	FinishReason string         `json:"finish_reason"`
}

// StreamDeltaObj represents the nested "delta" object in the JSON
type StreamDeltaObj struct {
	Content string `json:"content"`
}
