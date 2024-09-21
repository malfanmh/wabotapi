package usecase

import (
	"context"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"github.com/pkg/errors"
	"strings"
)

func (uc *useCase) sendMessageBySlug(ctx context.Context, clientID int64, WAPhoneID, WAID, slug string, access model.Access, param any) error {
	msg, err := uc.repo.GetMessageBySlug(ctx, clientID, slug)
	if err != nil {
		return err
	}
	jsonMessage, errM := uc.generateMessageBody(ctx, clientID, msg.ID, access)
	if errM != nil {
		return errM
	}
	fmt.Println(jsonMessage)
	jsonMessage, errT := uc.renderTemplate(jsonMessage, param)
	if errT != nil {
		fmt.Println("renderTemplate err:", errT)
	}

	result, errS := uc.wa.Send(ctx, WAPhoneID, WAID, msg.Type.ToWaType(), jsonMessage)
	if errS != nil {
		err = errors.WithStack(errS)
	}
	fmt.Println(result)
	return nil
}

func (uc *useCase) generateMessageBody(ctx context.Context, clientID, messageID int64, access model.Access) (jsonMessage string, err error) {
	message, err := uc.repo.GetMessage(ctx, clientID, messageID)
	if err != nil {
		return
	}
	jsonMessage = `{`
	switch message.Type {
	case model.MessageFlowList:
		jsonMessage += `"type":"list",`
		//TODO handle header type image
		if message.HeaderText != "" {
			jsonMessage += fmt.Sprintf(`"header":{"type":"text", "text":"%s"},`, message.HeaderText)
		}

		if message.BodyText != "" {
			jsonMessage += fmt.Sprintf(`"body":{"text":"%s"},`, message.BodyText)
		}
		if message.FooterText != "" {
			jsonMessage += fmt.Sprintf(`"footer":{"text":"%s"},`, message.FooterText)
		}

		actions, errA := uc.repo.GetMessageAction(ctx, messageID, access)
		if errA != nil {
			err = errA
			return
		}
		var actionJson []string
		for _, action := range actions {
			if action.Desc != "" {
				actionJson = append(actionJson, fmt.Sprintf(`{"id":"%s","title":"%s","description": "%s"}`, action.Slug, action.Title, action.Desc))
			} else {
				actionJson = append(actionJson, fmt.Sprintf(`{"id":"%s","title":"%s"}`, action.Slug, action.Title))
			}
		}
		jsonMessage += fmt.Sprintf(`"action":{"button":"%s","sections":[{"rows":[%s]}]}`, message.Button, strings.Join(actionJson, ","))

	case model.MessageFlowButton:
		jsonMessage += `"type":"button",`
		if message.HeaderText != "" {
			jsonMessage += fmt.Sprintf(`"header":{"type":"text", "text","%s"},`, message.HeaderText)
		}

		if message.BodyText != "" {
			jsonMessage += fmt.Sprintf(`"body":{"text":"%s"},`, message.BodyText)
		}
		if message.FooterText != "" {
			jsonMessage += fmt.Sprintf(`"footer":{"text":"%s"},`, message.FooterText)
		}

		actions, errA := uc.repo.GetMessageAction(ctx, messageID, access)
		if errA != nil {
			err = errA
			return
		}
		var actionJson []string
		for _, action := range actions {
			actionJson = append(actionJson, fmt.Sprintf(`{"type":"reply","reply":{"id":"%s","title":"%s"}}`, action.Slug, action.Title))
		}
		jsonMessage += fmt.Sprintf(`"action":{"buttons":[%s]}`, strings.Join(actionJson, ","))
	case model.MessageFlowText:
		if message.PreviewURL {
			jsonMessage += `"preview_url": true,`
		}

		if message.BodyText != "" {
			jsonMessage += fmt.Sprintf(`"body":"%s"`, message.BodyText)
		}
	case model.MessageFlowCTAURL:
		jsonMessage += `"type":"cta_url",`
		if message.HeaderText != "" {
			jsonMessage += fmt.Sprintf(`"header":{"type":"text", "text","%s"},`, message.HeaderText)
		}

		if message.BodyText != "" {
			jsonMessage += fmt.Sprintf(`"body":{"text":"%s"},`, message.BodyText)
		}
		if message.FooterText != "" {
			jsonMessage += fmt.Sprintf(`"footer":{"text":"%s"},`, message.FooterText)
		}

		actions, errA := uc.repo.GetMessageAction(ctx, messageID, access)
		if errA != nil {
			err = errA
			return
		}
		var actionJson string
		for _, action := range actions {
			actionJson = fmt.Sprintf(`{"display_text":"%s","url":"%s"}`, action.Title, action.Desc)
		}
		jsonMessage += fmt.Sprintf(`"action":{"name":"cta_url","parameters": %s }`, actionJson)
	default:
		err = fmt.Errorf("unknown type: %s", message.Type)
	}
	jsonMessage += "}"
	return
}
