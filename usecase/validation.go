package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"net/mail"
	"strconv"
	"strings"
	"time"
)

func (uc *useCase) validateRegistrationFlow(ctx context.Context, client model.Client, session model.Session, input string, customer model.Customer, flow model.MessageFlow) (ses model.Session) {
	// PENDAFTARAN
	isInputRevalidation := len(strings.Split(session.Slug, ":")) == 3
	if isInputRevalidation {
		flow.IsReValidate = false
	}

	slug := session.Slug
	if isInputRevalidation || flow.IsReValidate {
		slug = flow.Slug
	}
	fmt.Println("isInputRevalidation", isInputRevalidation, flow.IsReValidate, session.Slug, slug)
	keys := strings.Split(slug, ":")
	if len(keys) == 2 {
		if customer.WAID == "" {
			customer.WAID = session.WAID
			customer.ClientID = session.ClientID
			_ = uc.repo.InsertCustomer(ctx, customer)
		}

		if keys[0] == "menu-pendaftaran" {
			valid := true
			text := ""
			switch keys[1] {
			case "Nama":
				if flow.IsReValidate {
					msg, _ := uc.repo.GetMessage(ctx, client.ID, flow.MessageID)
					_, _ = uc.wa.SendText(ctx, client.WAPhoneID, customer.WAID, msg.BodyText, nil)
					session.Slug = flow.Slug + ":re"
					session.Seq = flow.Seq
					ses = session
					return
				}
				validateInput := strings.ReplaceAll(input, " ", "")
				if !model.IsLetter(validateInput) {
					text = fmt.Sprintf(`Nama yang anda inputkan (%s) salah, silahkan masukan nama yg benar\n\ncontoh: Ahmad Sayuri`, input)
					valid = false
				}
				customer.FullName = sql.Null[string]{
					V:     input,
					Valid: true,
				}
			case "TglLahir":
				if flow.IsReValidate {
					msg, _ := uc.repo.GetMessage(ctx, client.ID, flow.MessageID)
					_, _ = uc.wa.SendText(ctx, client.WAPhoneID, customer.WAID, msg.BodyText, nil)
					session.Slug = flow.Slug + ":re"
					session.Seq = flow.Seq
					ses = session
					return
				}

				t, err := time.Parse("02/01/2006", input)
				if err != nil {
					text = fmt.Sprintf(`Tanggal Lahir yang anda masukan (%s) salah, silahkan masukan dengan format\n\ncontoh: 17/08/1945`, input)
					valid = false
					if customer.BirthDate.Valid {
						t, _ = time.Parse(time.RFC3339, customer.BirthDate.V)
						customer.BirthDate = sql.Null[string]{
							V:     t.Format(time.DateOnly),
							Valid: true,
						}
					}

				} else {
					customer.BirthDate = sql.Null[string]{
						V:     t.Format(time.DateOnly),
						Valid: true,
					}
				}
			case "Alamat":
				if flow.IsReValidate {
					msg, _ := uc.repo.GetMessage(ctx, client.ID, flow.MessageID)
					_, _ = uc.wa.SendText(ctx, client.WAPhoneID, customer.WAID, msg.BodyText, nil)
					session.Slug = flow.Slug + ":re"
					session.Seq = flow.Seq
					ses = session
					return
				}

				if input == "" {
					text = `Masukan Alamat Sesuai KTP`
					valid = false
				}
				customer.Address = sql.Null[string]{
					V:     input,
					Valid: true,
				}
			case "JenisKelamin":
				if flow.IsReValidate {
					msg, _ := uc.repo.GetMessage(ctx, client.ID, flow.MessageID)
					_, _ = uc.wa.SendText(ctx, client.WAPhoneID, customer.WAID, msg.BodyText, nil)
					session.Slug = flow.Slug + ":re"
					session.Seq = flow.Seq
					ses = session
					return
				}

				switch strings.ToLower(input) {
				case "l", "p":
				default:
					text = fmt.Sprintf(`Jenis kelamin yang anda masukan (%s) salah salah, silakan masukan Jenis Kelamin yang benar:\n\ncontoh: \n*L* untuk laki-laki \n*P* untuk perempuan`, input)
					valid = false
				}
				customer.Gender = sql.Null[string]{
					V:     strings.ToUpper(input),
					Valid: true,
				}
			case "Email":
				if flow.IsReValidate {
					msg, _ := uc.repo.GetMessage(ctx, client.ID, flow.MessageID)
					_, _ = uc.wa.SendText(ctx, client.WAPhoneID, customer.WAID, msg.BodyText, nil)
					session.Slug = flow.Slug + ":re"
					session.Seq = flow.Seq
					ses = session
					return
				}

				if _, err := mail.ParseAddress(input); err != nil {
					text = fmt.Sprintf(`Email yang anda masukan (%s) salah, gunakan format email yang sesuai \n\ncontoh: ahmad.sayuri@muhammadiyah.com`, input)
					valid = false
				} else {
					customer.Email = sql.Null[string]{
						V:     input,
						Valid: true,
					}
				}

			case "Ktp":
				if flow.IsReValidate {
					msg, _ := uc.repo.GetMessage(ctx, client.ID, flow.MessageID)
					_, _ = uc.wa.SendText(ctx, client.WAPhoneID, customer.WAID, msg.BodyText, nil)
					session.Slug = flow.Slug + ":re"
					session.Seq = flow.Seq
					ses = session
					return
				}

				if len(input) != 16 {
					valid = false
					text = fmt.Sprintf(`Nomor KTP yang anda masukan (%s) salah, silahkan masukan 16 digit nomor KTP yang benar`, input)
				}

				if _, errNumber := strconv.Atoi(input); errNumber != nil {
					valid = false
					text = fmt.Sprintf(`Nomor KTP yang anda masukan (%s) salah, silahkan masukan 16 digit nomor KTP yang benar`, input)
				}

				customer.IdentityNumber = sql.Null[string]{
					V:     input,
					Valid: true,
				}
			}

			if !valid {
				result, err := uc.wa.SendText(ctx, client.WAPhoneID, customer.WAID, text, nil)
				fmt.Println(result, err)
				ses = session
				return
			}

			if err := uc.repo.UpdateCustomer(ctx, customer); err != nil {
				fmt.Println("UpdateCustomer err:", err)
			}
			if isInputRevalidation || flow.IsReValidate {
				flow, err := uc.repo.GetMessageFlowBySlug(ctx, client.ID, "menu-pendaftaran-preview")
				if err != nil {
					fmt.Println("get next flow err:", err)
				}
				ses, err = uc.regularFlow(ctx, client, flow, session, customer)
				if err != nil {
					fmt.Println("regularFlow err:", err)
				}
			} else {
				flow, err := uc.repo.GetNextFlow(ctx, client.ID, session.Access, keys[0], session.Seq)
				if err != nil {
					fmt.Println("get next flow err:", err)
				}
				ses, err = uc.regularFlow(ctx, client, flow, session, customer)
				if err != nil {
					fmt.Println("regularFlow err:", err)
				}
			}

		}
	}
	return
}

func (uc *useCase) validateActivationFlow(ctx context.Context, client model.Client, session model.Session, input string, customer model.Customer, flow model.MessageFlow) (ses model.Session) {
	if !model.IsAlphanumeric(input) {
		result, err := uc.wa.SendText(ctx, client.WAPhoneID, customer.WAID, fmt.Sprintf("Nomor KTA (%s) yang anda masukan salah \nmasukan nomor KTA yang benar.", input), nil)
		fmt.Println(result, err)
		ses = session
		return
	}
	customer.Status = sql.Null[model.Access]{
		V:     model.AccessActivated,
		Valid: true,
	}
	if err := uc.repo.UpdateCustomer(ctx, customer); err != nil {
		fmt.Println("UpdateCustomer err:", err)
	}
	session.Access = model.AccessActivated
	flow, err := uc.repo.GetMessageFlowBySlug(ctx, client.ID, "menu-aktivasi-done")
	if err != nil {
		fmt.Println("get next flow err:", err)
	}
	ses, err = uc.regularFlow(ctx, client, flow, session, customer)
	if err != nil {
		fmt.Println("regularFlow err:", err)
	}
	return
}
