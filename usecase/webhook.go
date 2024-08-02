package usecase

import (
	"context"
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
					return fmt.Errorf("unmarshal message err: %v", err)
				}

				switch msg.Type {
				case model.WAMessageTypeText:
					keyword, valid := msg.Text.Match([]string{"halo admin"}...)
					if !valid {
						fmt.Println("invalid keyword", msg.Text)
						return nil
					}

					if err = uc.staticFlow(ctx, keyword, msg.Text.Body, contact); err != nil {
						return err
					}
				case model.WAMessageTypeButton:
					if err = uc.staticFlow(ctx, msg.Button.Text, "", contact); err != nil {
						return err
					}
				case model.WAMessageTypeInteractive:
					if err = uc.staticFlow(ctx, msg.Interactive.ListReplay.ID, "", contact); err != nil {
						return err
					}
				default:
					fmt.Println("unknown type:", msg.Type)
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

const (
	salam           = `{ "body":"Assalamualaikum %s Selamat datang di Layanan WA Muhammadiyah. \n\nSebelum melanjutkan, harap masukan terlebih dahulu No KTA dan Nama Anggota Anda dengan format: \n\nNo. KTA#Nama Anggota\n123456789#Ahmad Sayuri"}`
	salamRegistered = `{"type":"list","body":{"text":"Assalamualaikum %s, Selamat datang di Layanan WA Muhammadiyah. \n\nSilahkan pilih layanan di bawah ini sesuai kebutuhan Anda."},"action":{"button":"Pilih Layanan","sections":[{"rows":[{"id":"aktivasi_anggota","title":"Aktivasi Anggota"},{"id":"pembayaran_iuran_anggota","title":"Pembayaran Iuran Anggota"},{"id":"pembayaran_zis","title":"Pembayaran ZIS","description":"(Zakat, Infaq, Sedekah)"},{"id":"etalase_produk","title":"Etalase Produk"},{"id":"informasi_umum","title":"Informasi Umum"}]}]},"footer":{"text":"%s"}}`
	informasiUmum   = `{"type":"list","body":{"text":"Informasi apa yang ingin Anda ketahui mengenai Koperasi Muhammadiyah?\n"},"action":{"button":"Jenis Informasi","sections":[{"rows":[{"id":"info_profil_koperasi","title":"Profil Koperasi","description":"Informasi Profil Koperasi Muhammadia"},{"id":"info_aktivasi_anggota","title":"Aktivasi Anggota"},{"id":"info_keanggotaan","title":"Informasi Keanggotaan"},{"id":"info_layanan_produk","title":"Informasi Layanan Produk"},{"id":"etalase_produk","title":"Etalase Produk"}]}]},"footer":{"text":"%s"}}`
	informasiProfil = `{"body":"Anda akan diarahkan menuju website koperasi Muhammadiyah \n\nhttps://muhammadiyah.or.id","preview_url":true}`
)

func (uc *useCase) staticFlow(ctx context.Context, keyword, payload string, contact model.WAContact) error {
	senderNumberID := "385924484596973"
	keyword = strings.ToLower(keyword)

	member, exists := uc.members.Get(contact.WaID)
	var msgType, jsonBody string
	switch keyword {
	case "halo admin":
		if !exists {
			name := func() string {
				if contact.Profile.Name.IsCleanLetter() {
					return contact.Profile.Name.String()
				}
				return ""
			}()
			msgType = model.WAMessageTypeText.String()
			jsonBody = fmt.Sprintf(salam, name)
		} else {
			msgType = model.WAMessageTypeInteractive.String()
			jsonBody = fmt.Sprintf(salamRegistered, member.Name, member.ID)
		}
	case model.KeywordInputKTA:
		_, _ = uc.wa.Send(ctx, senderNumberID, contact.WaID, model.WAMessageTypeText.String(),
			`{"body":"Mohon tunggu...."}`)
		data := strings.Split(payload, "#")
		uc.members.Set(contact.WaID, model.Member{
			ID:      data[0],
			Name:    data[1],
			Contact: contact,
		})
		_, _ = uc.wa.Send(ctx, senderNumberID, contact.WaID, model.WAMessageTypeText.String(),
			`{"body":"Aktivasi Anggota Berhasil...."}`)

		msgType = model.WAMessageTypeInteractive.String()
		jsonBody = fmt.Sprintf(salamRegistered, data[1], data[0])
	case "informasi_umum":
		msgType = model.WAMessageTypeInteractive.String()
		jsonBody = fmt.Sprintf(informasiUmum, member.ID)
	case "info_profil_koperasi":
		msgType = model.WAMessageTypeText.String()
		jsonBody = fmt.Sprintf(informasiProfil)
	default:
		fmt.Println("unknown keyword:", keyword)
		return nil
	}
	result, err := uc.wa.Send(ctx, senderNumberID, contact.WaID, msgType, jsonBody)
	fmt.Println("send response: ", result, err)
	return err
}
