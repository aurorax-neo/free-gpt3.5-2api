package respmodel

type ChatResp35 struct {
	Message struct {
		Id     string `json:"id"`
		Author struct {
			Role     string      `json:"role"`
			Name     interface{} `json:"name"`
			Metadata struct {
			} `json:"metadata"`
		} `json:"author"`
		CreateTime float64     `json:"create_time"`
		UpdateTime interface{} `json:"update_time"`
		Content    struct {
			ContentType string   `json:"content_type"`
			Parts       []string `json:"parts"`
		} `json:"content"`
		Status   string      `json:"status"`
		EndTurn  interface{} `json:"end_turn"`
		Weight   float64     `json:"weight"`
		Metadata struct {
			Citations        []interface{} `json:"citations"`
			GizmoId          interface{}   `json:"gizmo_id"`
			MessageType      string        `json:"message_type"`
			ModelSlug        string        `json:"model_slug"`
			DefaultModelSlug string        `json:"default_model_slug"`
			Pad              string        `json:"pad"`
			ParentId         string        `json:"parent_id"`
		} `json:"metadata"`
		Recipient string `json:"recipient"`
	} `json:"message"`
	ConversationId string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
}
