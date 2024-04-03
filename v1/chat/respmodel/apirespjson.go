package respmodel

type JsonObj struct {
	ID      string          `json:"id"`
	Object  string          `json:"object"`
	Created int64           `json:"created"`
	Model   string          `json:"model"`
	Usage   JsonUsageObj    `json:"usage"`
	Choices []JsonChoiceObj `json:"choices"`
}

type JsonMessageObj struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type JsonChoiceObj struct {
	Message      JsonMessageObj `json:"message"`
	FinishReason string         `json:"finish_reason"`
	Index        int            `json:"index"`
}

type JsonUsageObj struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
