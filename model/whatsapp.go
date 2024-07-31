package model

type (
	WAContactProfile struct {
		Name string `json:"name"`
	}
	WAContact struct {
		WaID    string           `json:"wa_id"`
		Profile WAContactProfile `json:"profile"`
	}
	WAMessageText struct {
		Body string `json:"body"`
	}
	WAMessageButton struct {
		Text    string `json:"text"`
		Payload string `json:"payload"`
	}
	WAMessage struct {
		From      string          `json:"from"`
		ID        string          `json:"id"`
		Timestamp string          `json:"timestamp"`
		Text      WAMessageText   `json:"text"`
		Button    WAMessageButton `json:"button"`
		Type      string          `json:"type"`
	}

	WATemplate struct {
		ID         string                   `json:"id"`
		Name       string                   `json:"name"`
		Language   string                   `json:"language"`
		Status     string                   `json:"status"`
		Category   string                   `json:"category"`
		Components []map[string]interface{} `json:"components"`
	}
)
