package reqmodel

type ApiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ApiReq struct {
	Messages    []ApiMessage `json:"messages"`
	Model       string       `json:"model"`
	Stream      bool         `json:"stream"`
	PluginIds   []string     `json:"plugin_ids"`
	NewMessages string
}
