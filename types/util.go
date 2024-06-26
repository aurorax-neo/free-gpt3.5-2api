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
		"gpt-4o":                 "gpt-4o",
		"auto":                   "auto",
		"gpt-4o-av":              "gpt-4o-av",
	}
	if v, ok := modelMapping[model]; ok {
		return v
	}
	return ""
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
		ForceUseSse:                true,
		FaceUseSse:                 false,
		ConversationMode: ChatConversationMode{
			Kind: "primary_assistant",
		},
		WebsocketRequestId: uuid.New().String(),
	}
	return chatReq
}
