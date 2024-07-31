package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"github.com/tidwall/gjson"
	"strings"
)

func (uc *useCase) VerifyWebhook(ctx context.Context, hash, token string) (bool, error) {
	client, err := uc.repo.GetClientByHash(ctx, hash)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	fmt.Println(client.Token, token)
	return token == client.Token, nil
}

func (uc *useCase) Webhook(ctx context.Context, hash string, payload []byte) error {
	data := gjson.ParseBytes(payload)
	client, err := uc.repo.GetClientByHash(ctx, hash)
	if err != nil {
		fmt.Println("GetClientByHash err:", err)
		return err
	}
	if client.Hash != hash {
		fmt.Println("client hash does not match")
		return nil
	}
	entryChanges := data.Get("entry.0.changes.0")
	field := entryChanges.Get("field").String()
	switch field {
	case "messages":
		phoneNumber := entryChanges.Get("value.contacts.0.wa_id").String()
		if entryChanges.Get("value.messages").Exists() {
			strMessages := entryChanges.Get("value.messages").String()
			fmt.Println(strMessages)
			var messages []model.WAMessage
			if err = json.Unmarshal([]byte(strMessages), &messages); err != nil {
				fmt.Println("unmarshal messages err:", err)
				return nil
			}
			for _, msg := range messages {
				switch msg.Type {
				case "text":
					if err = uc.staticFlow(ctx, msg.Text.Body, phoneNumber); err != nil {
						return err
					}
				case "button":
					if err = uc.staticFlow(ctx, msg.Button.Text, phoneNumber); err != nil {
						return err
					}
				case "link":
				default:
					fmt.Println("unknown type:", msg.Type)
				}
			}
		}
		// TODO handle message status
	default:
		fmt.Println("unknown field:", field)
		return nil
	}
	return nil
}

func (uc *useCase) staticFlow(ctx context.Context, keyword, to string) error {
	senderNumberID := "385924484596973"
	keyword = strings.ToLower(keyword)
	switch keyword {
	case "hi":
		result, err := uc.wa.SendTemplate(ctx,
			senderNumberID,
			to,
			model.WATemplate{
				Name:     "hello_world",
				Language: "en_US",
			}, nil)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("send ", result)
		result, err = uc.wa.SendTemplate(ctx,
			senderNumberID,
			to,
			model.WATemplate{
				Name:     "qna_template_test_1",
				Language: "en",
			}, nil)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("send ", result)
	case "produk a":
		result, err := uc.wa.SendText(ctx, senderNumberID, to, "terimakasih sudah memilih produk A", nil)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("send ", result)
	case "produk b":
		result, err := uc.wa.SendText(ctx, senderNumberID, to, "terimakasih sudah memilih produk B", nil)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("send ", result)
	}
	return nil
}
