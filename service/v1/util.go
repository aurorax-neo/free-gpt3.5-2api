package v1

import (
	"free-gpt3.5-2api/service/v1Chat/reqModel"
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
	if model == "" {
		return "text-davinci-002-render-sha"
	}
	if v, ok := modelMapping[model]; ok {
		return v
	}
	return "text-davinci-002-render-sha"
}

func GenerateID(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	id := "chatcmpl-"
	for i := 0; i < length; i++ {
		id += string(charset[rand.Intn(len(charset))])
	}
	return id
}

func ApiReq2ChatReq35(apiReq *reqModel.ApiReq) (chatReq *reqModel.ChatReq35) {
	messages := make([]reqModel.ChatMessages, 0)
	for _, apiMessage := range apiReq.Messages {
		chatMessage := reqModel.ChatMessages{
			Author: reqModel.ChatAuthor{
				Role: apiMessage.Role,
			},
			Content: reqModel.ChatContent{
				ContentType: "text",
				Parts:       []string{apiMessage.Content},
			},
		}
		messages = append(messages, chatMessage)
	}

	chatReq = &reqModel.ChatReq35{
		Action:                     "next",
		Messages:                   messages,
		ParentMessageId:            uuid.New().String(),
		Model:                      MappingModel(apiReq.Model),
		TimeZoneOffsetMin:          -180,
		Suggestions:                make([]string, 0),
		HistoryAndTrainingDisabled: true,
		ConversationMode: reqModel.ChatConversationMode{
			Kind: "primary_assistant",
		},
		WebsocketRequestId: uuid.New().String(),
	}
	return chatReq
}
