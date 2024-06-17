package types

import "time"

type ChatResp struct {
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
			ContentType string        `json:"content_type"`
			Parts       []interface{} `json:"parts"`
			Language    string        `json:"language"`
			Text        string        `json:"text"`
		} `json:"content"`
		Status    string      `json:"status"`
		EndTurn   interface{} `json:"end_turn"`
		Weight    float64     `json:"weight"`
		Metadata  Metadata    `json:"metadata"`
		Recipient string      `json:"recipient"`
	} `json:"message"`
	ConversationId string      `json:"conversation_id"`
	Error          interface{} `json:"error"`
	// 审核
	Type               string `json:"types"`
	MessageId          string `json:"message_id"`
	IsCompletion       bool   `json:"is_completion"`
	ModerationResponse struct {
		Flagged      bool          `json:"flagged"`
		Disclaimers  []interface{} `json:"disclaimers"`
		Blocked      bool          `json:"blocked"`
		ModerationId string        `json:"moderation_id"`
	} `json:"moderation_response"`
}

type DalleContent struct {
	AssetPointer string `json:"asset_pointer"`
	Metadata     struct {
		Dalle struct {
			Prompt string `json:"prompt"`
		} `json:"dalle"`
	} `json:"metadata"`
}

type Metadata struct {
	Timestamp     string         `json:"timestamp_"`
	Citations     []Citation     `json:"citations,omitempty"`
	MessageType   string         `json:"message_type"`
	FinishDetails *FinishDetails `json:"finish_details"`
	ModelSlug     string         `json:"model_slug"`

	GizmoId           interface{} `json:"gizmo_id"`
	DefaultModelSlug  string      `json:"default_model_slug"`
	Pad               string      `json:"pad"`
	ParentId          string      `json:"parent_id"`
	ModelSwitcherDeny []struct {
		Slug        string `json:"slug"`
		Context     string `json:"context"`
		Reason      string `json:"reason"`
		Description string `json:"description"`
	} `json:"model_switcher_deny"`
}

type Citation struct {
	Metadata CitaMeta `json:"metadata"`
	StartIx  int      `json:"start_ix"`
	EndIx    int      `json:"end_ix"`
}

type CitaMeta struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type FinishDetails struct {
	Type string `json:"types"`
	Stop string `json:"stop"`
}

type ChatRespInfo struct {
	Type         string    `json:"types"`
	CallToAction string    `json:"call_to_action"`
	ResetsAfter  time.Time `json:"resets_after"`
	LimitDetails struct {
		Type                  string `json:"types"`
		ModelSlug             string `json:"model_slug"`
		UsingDefaultModelSlug string `json:"using_default_model_slug"`
		NextModelSlug         string `json:"next_model_slug"`
		ModelLimitName        string `json:"model_limit_name"`
	} `json:"limit_details"`
	DisplayDescription struct {
		Type                string      `json:"types"`
		Description         string      `json:"description"`
		MarkdownDescription interface{} `json:"markdown_description"`
	} `json:"display_description"`
	ConversationId string `json:"conversation_id"`
}

type GenericResponseLine struct {
	Line  string `json:"line"`
	Error string `json:"error"`
}

type StringStruct struct {
	Text string `json:"text"`
}
