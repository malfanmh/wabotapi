package model

import (
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

func isLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func (pn WAContactProfileName) IsCleanLetter() bool {
	return isLetter(pn.String())
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
		Type       string                  `json:"type"`
		ListReplay WAInteractiveListReplay `json:"list_reply"`
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

const KeywordInputKTA = "input_kta"

func (t WAMessageText) Match(keywords ...string) (keyword string, valid bool) {
	if t.IsInputKTA() {
		return KeywordInputKTA, true
	}
	for _, v := range keywords {
		if strings.ToLower(v) == t.Body {
			return t.Body, true
		}
	}
	return "", false
}

func (wt WAMessageType) String() string {
	return string(wt)
}

type Member struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Contact WAContact `json:"contact"`
}
