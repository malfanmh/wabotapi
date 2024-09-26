package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"log"
	"math/rand"
)

func (uc *useCase) checkout(ctx context.Context, client model.Client, session model.Session, input string, customer model.Customer, flow model.MessageFlow) (ses model.Session) {
	randomNumber := func(min, max int) int64 {
		return int64(rand.Intn(max-min) + min)
	}
	// TODO get items
	amount := decimal.NewFromInt(randomNumber(15000, 50000))
	product, err := uc.repo.GetProductBySlug(ctx, client.ID, input)
	if err != nil {
		log.Println(err)
	}
	amount = product.Price

	orderID, err := uc.repo.CreatePayment(ctx, model.Payment{
		CustomerID: customer.ID,
		ClientID:   client.ID,
		RefID: sql.Null[string]{
			V:     "",
			Valid: true,
		},
		PaymentProvider: sql.Null[string]{
			V:     "",
			Valid: true,
		},
		PaymentItem: sql.Null[string]{
			V:     product.Slug,
			Valid: input != "",
		},
		Status: sql.Null[string]{
			V:     "PENDING",
			Valid: true,
		},
		Amount: decimal.NewNullDecimal(amount),
		Fee:    decimal.NewNullDecimal(decimal.NewFromInt(0)),
	})
	if err != nil {
		fmt.Println("CreatePayment err:", err)
		return session
	}

	paymentLink, err := uc.payment.GetPaymentLink(ctx, client, model.PaymentLink{
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

	//fmt.Println(paymentLink.Redirecturl, err)
	ses, err = uc.regularFlow(ctx, client, flow, session, model.MessageMetadata{
		ID:          fmt.Sprint(customer.ID),
		Name:        customer.FullName.V,
		ProductName: product.Name,
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
		fmt.Println("GetPaymentCustomer not found, err:", err)
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
		err = uc.sendMessageBySlug(ctx, ct.ClientID, ct.WAPhoneID, ct.WAID, "checkout-payment-success", model.AccessActivated,
			model.MessageMetadata{
				Invoice: fmt.Sprintf("INV%d", ct.ID),
				ID:      fmt.Sprint(ct.CustomerID),
				Name:    ct.FullName,
			})
		fmt.Println("sendMessageBySlug", err)
	case "CANCELLED", "FAILURE":
		// TODO send message , checkout-payment-failed
		err = uc.sendMessageBySlug(ctx, ct.ClientID, ct.WAPhoneID, ct.WAID, "checkout-payment-failed", model.AccessActivated, nil)
		return nil
	default:
		fmt.Println("Unhandled Payment status:", callback.Result.Payment.Status)
	}
	amount := decimal.NewFromFloat(callback.Result.Payment.Amount)
	data := model.Payment{
		ID:              ct.ID,
		RefID:           toNullString(callback.Order.Reference),
		CustomerID:      ct.CustomerID,
		PaymentType:     toNullString(callback.SourceOfFunds.Type),
		PaymentProvider: toNullString("finpay"),
		PaymentCode:     toNullString(callback.SourceOfFunds.PaymentCode),
		PaymentRefID:    toNullString(callback.Result.Payment.Reference),
		Status:          toNullString(callback.Result.Payment.Status),
		Amount:          decimal.NewNullDecimal(amount),
		Fee:             decimal.NewNullDecimal(decimal.Zero),
	}

	if err = uc.repo.UpdatePayment(ctx, data); err != nil {
		return err
	}
	return nil
}

func (uc *useCase) ExpiryLink(ctx context.Context) error {
	payments, err := uc.repo.GetExpiredPayment(ctx)
	if err != nil {
		return err
	}
	for _, payment := range payments {
		errU := uc.repo.UpdatePayment(ctx, model.Payment{
			ID: payment.ID,
			Status: sql.Null[string]{
				V:     "EXPIRED",
				Valid: true,
			},
		})
		if errU != nil {
			err = errors.WithStack(errU)
		}
	}
	if err != nil {
		return err
	}
	return nil
}
