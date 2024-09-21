package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"github.com/pkg/errors"
	"strings"
)

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
		session.Access = customer.Access.V
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
		switch {
		case flow.ValidateInput:
			if flow.IsInput {
				switch session.Access {
				case model.AccessNew:
					session = uc.validateRegistrationFlow(ctx, client, session, payload, customer, flow)
				case model.AccessRegistered:
					session = uc.validateActivationFlow(ctx, client, session, payload, customer, flow)
				}
			} else {
				// TODO optimize this flow
				switch keyword {
				case "menu-pendaftaran-preview-ya":
					fmt.Println("menu-pendaftaran-preview-ya", model.AccessRegistered)
					_ = uc.repo.UpdateCustomer(ctx, model.Customer{WAID: contact.WaID, ClientID: client.ID, Access: sql.Null[model.Access]{
						V:     model.AccessRegistered,
						Valid: true,
					}})
					session.Access = model.AccessRegistered
				case model.InputTextKTA:
					fmt.Println(model.InputTextKTA, model.AccessActivated)
					_ = uc.repo.UpdateCustomer(ctx, model.Customer{WAID: contact.WaID, ClientID: client.ID, Access: sql.Null[model.Access]{
						V:     model.AccessActivated,
						Valid: true,
					}})
					session.Access = model.AccessActivated
				}
				session, _ = uc.regularFlow(ctx, client, flow, session, md)
			}
		case flow.Checkout:
			session = uc.checkout(ctx, client, session, keyword, customer, flow)
		default:
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

// slug structure : keyword:entity:revalidate
// i.e : menu-pendaftaran:Nama:re

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
		} else if len(keys) == 1 {
			if keys[0] == "menu-aktivasi" {
				flows, err = uc.repo.GetMessageFlow(ctx, session.ClientID, session.Access, keys[0], session.Seq, 1)
				if err != nil {
					return nil, err
				}
				flows[0].IsInput = true
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
