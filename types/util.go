package types

import (
	"github.com/google/uuid"
	"math/rand"
)

func MappingModel(model string) string {
	var modelMapping = map[string]string{
		"gpt-3.5-turbo":          "text-davinci-002-render-sha",
		"gpt-3.5-turbo-16k":      "text-davinci-002-render-sha",
		"gpt-3.5-turbo-16k-0613": "text-davinci-002-render-sha",
		"gpt-3.5-turbo-0301":     "text-davinci-002-render-sha",
		"gpt-3.5-turbo-0613":     "text-davinci-002-render-sha",
		"gpt-3.5-turbo-1106":     "text-davinci-002-render-sha",
	}
	if v, ok := modelMapping[model]; ok {
		return v
	}
	return model
}

func GenerateID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := "chatcmpl-"
	for i := 0; i < length; i++ {
		id += string(charset[rand.Intn(len(charset))])
	}
	return id
}

func ApiReq2ChatReq35(apiReq *ApiReq) (chatReq *ChatReq) {
	messages := make([]ChatMessages, 0)
	for _, apiMessage := range apiReq.Messages {
		chatMessage := ChatMessages{
			Author: ChatAuthor{
				Role: apiMessage.Role,
			},
			Content: ChatContent{
				ContentType: "text",
				Parts:       []string{apiMessage.Content},
			},
		}
		messages = append(messages, chatMessage)
	}

	chatReq = &ChatReq{
		Action:                     "next",
		Messages:                   messages,
		ParentMessageId:            uuid.New().String(),
		Model:                      MappingModel(apiReq.Model),
		TimeZoneOffsetMin:          -180,
		Suggestions:                make([]string, 0),
		HistoryAndTrainingDisabled: true,
		ConversationMode: ChatConversationMode{
			Kind: "primary_assistant",
		},
		WebsocketRequestId: uuid.New().String(),
	}
	return chatReq
}
