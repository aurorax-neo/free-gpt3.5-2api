package types

type ChatAuthor struct {
	Role string `json:"role"`
}

type ChatContent struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type ChatMessages struct {
	Author  ChatAuthor  `json:"author"`
	Content ChatContent `json:"content"`
}

type ChatConversationMode struct {
	Kind string `json:"kind"`
}

type ChatReq struct {
	Action                     string               `json:"action"`
	Messages                   []ChatMessages       `json:"messages"`
	ParentMessageId            string               `json:"parent_message_id"`
	Model                      string               `json:"model"`
	TimeZoneOffsetMin          int                  `json:"timezone_offset_min"`
	Suggestions                []string             `json:"suggestions"`
	HistoryAndTrainingDisabled bool                 `json:"history_and_training_disabled"`
	ForceUseSse                bool                 `json:"force_use_sse"`
	FaceUseSse                 bool                 `json:"face_use_sse"`
	ConversationMode           ChatConversationMode `json:"conversation_mode"`
	WebsocketRequestId         string               `json:"websocket_request_id"`
}
