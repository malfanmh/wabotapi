package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"strings"
)

func (uc *useCase) VerifyWebhook(ctx context.Context, hash, token string) (bool, error) {
	client, err := uc.repo.GetClientByHash(ctx, hash)
	if err != nil {
		return false, err
	}
	return token == client.Token, nil
}

func (uc *useCase) Webhook(ctx context.Context, hash string, payload []byte) error {
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
				fmt.Println("message.String()", message.String())
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

func (uc *useCase) flow(ctx context.Context, client model.Client, contact model.WAContact, keyword, payload string) error {
	if keyword == "" {
		return nil
	}

	session, err := uc.repo.GetSession(ctx, client.ID, contact.WaID, keyword)
	if err != nil {
		return fmt.Errorf("get session err: %v", err)
	}
	customer, err := uc.repo.GetCustomerByWAID(ctx, client.ID, contact.WaID)
	if err == nil {
		session.Access = customer.Status.V
	}
	flows, err := uc.getFlow(ctx, session, keyword)
	if err != nil {
		return err
	}
	fmt.Println("flow", flows, err)
	for _, flow := range flows {
		md := model.MessageMetadata{
			Name: func() string {
				if customer.FullName.Valid {
					return customer.FullName.V
				}
				return contact.Profile.Name.String()
			}(),
		}
		fmt.Println("validate input", flow.ValidateInput, flow.IsInput)
		if flow.ValidateInput {
			if flow.IsInput {
				session = uc.validationFlow(ctx, client, session, payload, customer, flow)
			} else {
				// TODO optimize this flow
				switch keyword {
				case "menu-pendaftaran-preview-ya":
					fmt.Println("menu-pendaftaran-preview-ya", model.AccessRegistered)
					_ = uc.repo.UpdateCustomer(ctx, model.Customer{WAID: contact.WaID, ClientID: client.ID, Status: sql.Null[model.Access]{
						V:     model.AccessRegistered,
						Valid: true,
					}})
				case model.InputTextKTA:
					fmt.Println("menu-pendaftaran-preview-ya", model.AccessRegistered)
					_ = uc.repo.UpdateCustomer(ctx, model.Customer{WAID: contact.WaID, ClientID: client.ID, Status: sql.Null[model.Access]{
						V:     model.AccessActivated,
						Valid: true,
					}})
				}
				session, _ = uc.regularFlow(ctx, client, flow, session, md)
			}
		} else {
			session, _ = uc.regularFlow(ctx, client, flow, session, md)
		}
	}

	if err != nil {
		fmt.Println("send err:", err)
		return err
	}
	session.Input = payload
	err = uc.repo.UpdateSession(ctx, session)
	if err != nil {
		fmt.Println("UpdateSession err:", err)
		return err
	}
	return nil
}

func (uc *useCase) getFlow(ctx context.Context, session model.Session, input string) ([]model.MessageFlow, error) {
	if input == "menu-pendaftaran" {
		flows, err := uc.repo.GetMessageFlow(ctx, session.ClientID, session.Access, input, "0", 2)
		if err != nil {
			return nil, err
		}
		return flows, nil
	}

	keyword := input
	if keyword == session.Slug {
		keyword = session.Slug
	}
	flows, err := uc.repo.GetMessageFlow(ctx, session.ClientID, session.Access, keyword, "", 0)
	if err != nil {
		return nil, err
	}
	// is input text, check prev flow from session
	fmt.Println("ln flow", len(flows), err)

	if len(flows) == 0 {
		revalidate := false
		if session.Slug == "menu-pendaftaran-re-validate" {
			flow, err := uc.repo.GetMessageFlowBySlug(ctx, session.ClientID, input)
			fmt.Println("GetMessageFlowBySlug", err)
			session.Slug = flow.Slug
			session.Seq = flow.Seq
			revalidate = true
		}

		fmt.Println("session.Slug revalidate", session.Slug, revalidate)
		keys := strings.Split(session.Slug, ":")
		if len(keys) == 3 { // revalidation
			flows, err = uc.repo.GetMessageFlow(ctx, session.ClientID, session.Access, keys[0], session.Seq, 1)
			if err != nil {
				return nil, err
			}
			flows[0].IsInput = true
			flows[0].IsReValidate = true
		} else if len(keys) == 2 {
			if keys[0] == "menu-pendaftaran" {
				flows, err = uc.repo.GetMessageFlow(ctx, session.ClientID, session.Access, keys[0], session.Seq, 1)
				if err != nil {
					return nil, err
				}
				flows[0].IsInput = true
				flows[0].IsReValidate = revalidate
			}
		}
	}
	return flows, nil
}

func (uc *useCase) regularFlow(ctx context.Context, client model.Client, flow model.MessageFlow, session model.Session, args any) (ses model.Session, err error) {
	jsonMessage, errM := uc.generateMessageBody(ctx, client.ID, flow.MessageID, session.Access)
	if errM != nil {
		return session, errM
	}
	fmt.Println(jsonMessage)
	jsonMessage, errT := uc.renderTemplate(jsonMessage, args)
	if errT != nil {
		fmt.Println("renderTemplate err:", errT)
	}

	result, errS := uc.wa.Send(ctx, client.WAPhoneID, session.WAID, flow.Type.ToWaType(), jsonMessage)
	if errS != nil {
		err = errors.WithStack(errS)
	}
	session.Slug = flow.Slug
	if flow.ValidateInput {
		session.Seq = flow.Seq
	} else {
		session.Seq = "0"
	}
	fmt.Println(result)
	return session, nil
}
