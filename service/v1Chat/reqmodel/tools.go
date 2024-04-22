package reqmodel

import (
	"free-gpt3.5-2api/service/v1"
	"github.com/google/uuid"
)

func ApiReq2ChatReq35(apiReq *ApiReq) (chatReq *ChatReq35) {
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

	chatReq = &ChatReq35{
		Action:                     "next",
		Messages:                   messages,
		ParentMessageId:            uuid.New().String(),
		Model:                      v1.MappingModel(apiReq.Model),
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
