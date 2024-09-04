package model

import (
	"fmt"
	"strings"
	"unicode"
)

type WAMessageType string

const (
	WAMessageTypeText        WAMessageType = "text"
	WAMessageTypeButton      WAMessageType = "button"
	WAMessageTypeInteractive WAMessageType = "interactive"
)

type WAContactProfileName string

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func ParseWAMessageType(s string) (WAMessageType, error) {
	switch s {
	case "text":
		return WAMessageTypeText, nil
	case "list", "button":
		return WAMessageTypeInteractive, nil
	}
	return "", fmt.Errorf("invalid WA message type: %s", s)
}

func (pn WAContactProfileName) String() string {
	return string(pn)
}

type (
	WAContactProfile struct {
		Name WAContactProfileName `json:"name"`
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

	WAInteractiveListReplay struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	WAInteractive struct {
		Type         string                  `json:"type"`
		ListReplay   WAInteractiveListReplay `json:"list_reply"`
		ButtonReplay WAInteractiveListReplay `json:"button_reply"`
	}
	WAMessageContext struct {
		From string `json:"from"`
		To   string `json:"to"`
	}
	WAMessage struct {
		Context     WAMessageContext `json:"context"`
		From        string           `json:"from"`
		ID          string           `json:"id"`
		Timestamp   string           `json:"timestamp"`
		Text        WAMessageText    `json:"text"`
		Button      WAMessageButton  `json:"button"`
		Interactive WAInteractive    `json:"interactive"`
		Type        WAMessageType    `json:"type"`
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

func (t WAMessageText) IsInputKTA() bool {
	if strings.Contains(t.Body, "#") {
		input := strings.Split(t.Body, "#")
		if len(input) == 2 {
			return true
		}
	}
	return false
}

const (
	InputText    = "input_text"
	InputTextKTA = "input_text_kta"
)

func (t WAMessageText) Parse() (keyword, origin string, valid bool) {
	if t.IsInputKTA() {
		return InputTextKTA, t.Body, true
	}
	return t.Body, t.Body, t.Body != ""
}

func (wt WAMessageType) String() string {
	return string(wt)
}
