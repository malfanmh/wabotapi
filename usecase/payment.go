package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"github.com/shopspring/decimal"
	"math/rand"
)

func (uc *useCase) checkout(ctx context.Context, client model.Client, session model.Session, input string, customer model.Customer, flow model.MessageFlow) (ses model.Session) {
	randomNumber := func(min, max int) int64 {
		return int64(rand.Intn(max-min) + min)
	}
	// TODO get items
	amount := decimal.NewFromInt(randomNumber(15000, 50000))
	switch input {
	case "zis":
		amount = decimal.NewFromInt(randomNumber(50000, 100000))
	case "iuran":
		amount = decimal.NewFromInt(randomNumber(100000, 200000))
	}

	orderID, err := uc.repo.CreatePayment(ctx, model.Payment{
		CustomerID: customer.ID,
		PaymentItem: sql.Null[string]{
			V:     input,
			Valid: input != "",
		},
		Access: sql.Null[string]{
			V:     "PENDING",
			Valid: true,
		},
		Amount: decimal.NewNullDecimal(amount),
		Fee:    decimal.NullDecimal{},
	})
	if err != nil {
		fmt.Println("CreatePayment err:", err)
		return session
	}

	paymentLink, err := uc.payment.GetPaymentLink(ctx, model.PaymentLink{
		Customer: model.PaymentLinkCustomer{
			Email:       customer.Email.V,
			FirstName:   customer.FullName.V,
			LastName:    "-",
			MobilePhone: fmt.Sprintf("+%s", customer.WAID),
		},
		Order: model.PaymentLinkOrder{
			ID:          fmt.Sprint(orderID),
			Amount:      amount.String(),
			Description: input,
		},
	})
	if err != nil {
		fmt.Println("GetPaymentLink err:", err)
	}

	if err := uc.repo.UpdatePayment(ctx, model.Payment{
		ID: orderID,
		ExpiredAt: sql.Null[string]{
			V:     paymentLink.ExpiryLink,
			Valid: paymentLink.ExpiryLink != "",
		},
	}); err != nil {
		fmt.Println("UpdatePayment err:", err)
	}

	fmt.Println(paymentLink.Redirecturl, err)
	ses, err = uc.regularFlow(ctx, client, flow, session, model.MessageMetadata{
		CheckoutURL: paymentLink.Redirecturl,
		Amount:      model.FormatRP(amount.Truncate(2).InexactFloat64()),
		ExpiryDate:  model.FormatExpiryDate(paymentLink.ExpiryLink),
	})
	if err != nil {
		fmt.Println("regularFlow err:", err)
	}
	return
}

func (uc *useCase) PaymentCallback(ctx context.Context, callback model.FinpayCallback) error {
	ct, err := uc.repo.GetPaymentCustomer(ctx, callback.Order.ID)
	if err != nil {
		return err
	}
	var toNullString = func(s string) sql.Null[string] {
		return sql.Null[string]{
			V:     s,
			Valid: s != "",
		}
	}
	switch callback.Result.Payment.Status {
	case "PAID":
	case "CANCELLED", "FAILURE":

	}
	amount := decimal.NewFromFloat(callback.Result.Payment.Amount)
	data := model.Payment{
		ID:              ct.ID,
		RefID:           toNullString(callback.Order.Reference),
		CustomerID:      0,
		PaymentType:     toNullString(callback.SourceOfFunds.Type),
		PaymentProvider: toNullString("finpay"),
		PaymentCode:     toNullString(callback.SourceOfFunds.PaymentCode),
		PaymentRefID:    toNullString(callback.Result.Payment.Reference),
		Access:          toNullString(callback.Result.Payment.Status),
		Amount:          decimal.NewNullDecimal(amount),
		Fee:             decimal.NewNullDecimal(decimal.Zero),
	}

	if err = uc.repo.UpdatePayment(ctx, data); err != nil {
		return err
	}
	return nil
}
