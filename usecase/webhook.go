package usecase

import (
	"context"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"github.com/tidwall/gjson"
)

func (uc *useCase) VerifyWebhook(ctx context.Context, token string) (bool, error) {
	return token == uc.secret, nil
}

func (uc *useCase) Webhook(ctx context.Context, payload []byte) error {
	data := gjson.ParseBytes(payload)
	waPhoneID := data.Get("entry.0.changes.0.value.metadata.phone_number_id").String()
	client, err := uc.repo.GetClientByWAPhoneID(ctx, waPhoneID)
	if err != nil {
		return err
	}

	field := data.Get("entry.0.changes.0.field").String()
	switch field {
	case "messages":
		messageValue := data.Get("entry.0.changes.0.value")
		// get contact
		var contact model.WAContact
		if messageValue.Get("contacts").Exists() {
			strContact := messageValue.Get("contacts.0").String()
			if err = contact.UnmarshalJSON([]byte(strContact)); err != nil {
				return fmt.Errorf("unmarshal contact err: %v", err)
			}
		}

		if messageValue.Get("messages").Exists() {
			messages := messageValue.Get("messages").Array()
			for _, message := range messages {
				//fmt.Println("message.String()", message.String())
				var msg model.WAMessage
				if err := msg.UnmarshalJSON([]byte(message.String())); err != nil {
					fmt.Printf("unmarshal message err: %v\n", err)
					return fmt.Errorf("unmarshal message err: %v", err)
				}

				keyword := ""
				msgPayload := ""
				switch msg.Type {
				case model.WAMessageTypeText:
					value, origin, valid := msg.Text.Parse()
					if !valid {
						fmt.Println("invalid keyword", msg.Text)
						return nil
					}
					keyword = value
					msgPayload = origin
				case model.WAMessageTypeButton:
					keyword = msg.Button.Text
				case model.WAMessageTypeInteractive:
					if msg.Interactive.ListReplay.ID != "" {
						keyword = msg.Interactive.ListReplay.ID
					}
					if msg.Interactive.ButtonReplay.ID != "" {
						keyword = msg.Interactive.ButtonReplay.ID
					}
				default:
					fmt.Println("unknown type:", msg.Type)
					return nil
				}

				if err = uc.flow(ctx, client, contact, keyword, msgPayload); err != nil {
					return err
				}
			}
		} else {
			fmt.Println("message not found")
		}
		// TODO handle message statuses
	default:
		fmt.Println("unknown field:", field)
		return nil
	}
	return nil
}
