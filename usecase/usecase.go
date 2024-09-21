package usecase

import (
	"context"
	"github.com/malfanmh/wabotapi/model"
)

type Repository interface {
	GetClientByWAPhoneID(ctx context.Context, waPhoneID string) (result model.Client, err error)
	GetMessage(ctx context.Context, clientID, messageID int64) (msg model.Message, err error)
	GetMessageFlow(ctx context.Context, clientID int64, access model.Access, keyword string, seq string, limit int) (result []model.MessageFlow, err error)
	GetMessageFlowBySlug(ctx context.Context, clientID int64, slug string) (flow model.MessageFlow, err error)
	GetNextFlow(ctx context.Context, clientID int64, access model.Access, keyword string, seq string) (result model.MessageFlow, err error)
	GetMessageAction(ctx context.Context, messageID int64, access model.Access) (result []model.MessageAction, err error)

	GetCustomerByWAID(ctx context.Context, clientID int64, waid string) (result model.Customer, err error)
	InsertCustomer(ctx context.Context, customer model.Customer) (err error)
	UpdateCustomer(ctx context.Context, customer model.Customer) (err error)

	GetSession(ctx context.Context, clientID int64, waid, keyword string) (result model.Session, err error)
	UpdateSession(ctx context.Context, session model.Session) (err error)

	GetPaymentCustomer(ctx context.Context, refID string) (result model.PaymentCustomer, err error)
	CreatePayment(ctx context.Context, data model.Payment) (lastID int64, err error)
	UpdatePayment(ctx context.Context, data model.Payment) (err error)
}

type WhatsAppRepository interface {
	Send(ctx context.Context, from, to string, msgType model.WAMessageType, jsonBody string) (string, error)
	SendTemplate(ctx context.Context, from, to string, template model.WATemplate, params map[string]interface{}) (string, error)
	SendText(ctx context.Context, from, to string, text string, params map[string]interface{}) (string, error)
}

type Payment interface {
	GetPaymentLink(ctx context.Context, payload model.PaymentLink) (result model.PaymentLinkResponse, err error)
}

type useCase struct {
	repo    Repository
	wa      WhatsAppRepository
	payment Payment
	secret  string
}

func New(repository Repository, wa WhatsAppRepository, payment Payment, secret string) *useCase {
	return &useCase{
		repo:    repository,
		wa:      wa,
		payment: payment,
		secret:  secret,
	}
}
