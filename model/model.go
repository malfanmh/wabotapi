package model

import (
	"database/sql"
	"github.com/shopspring/decimal"
)

type Access int

const (
	AccessNew Access = iota
	AccessRegistered
	AccessActivated
)

func (a Access) String() string {
	return [...]string{"public", "registered", "activated"}[a]
}

func (a Access) Int() int {
	return int(a)
}

type Client struct {
	ID                int64  `db:"id"`
	Name              string `db:"name"`
	FinpayMerchantID  string `db:"finpay_merchant_id"`
	FinpaySecret      string `db:"finpay_secret"`
	FinpayCallbackURL string `db:"finpay_callback_url"`
	WAPhone           string `db:"wa_phone"`
	WAPhoneID         string `db:"wa_phone_id"`
}

type MessageFlowType string

const (
	MessageFlowText   MessageFlowType = "text"
	MessageFlowList   MessageFlowType = "list"
	MessageFlowButton MessageFlowType = "button"
	MessageFlowCTAURL MessageFlowType = "cta_url"
)

var mapType = map[MessageFlowType]WAMessageType{
	MessageFlowButton: WAMessageTypeInteractive,
	MessageFlowList:   WAMessageTypeInteractive,
	MessageFlowCTAURL: WAMessageTypeInteractive,
	MessageFlowText:   WAMessageTypeText,
}

func (f MessageFlowType) ToWaType() WAMessageType {
	return mapType[f]
}

type MessageFlow struct {
	Keyword       string          `db:"keyword"`
	MessageID     int64           `db:"message_id"`
	Access        Access          `db:"access"`
	Seq           string          `db:"seq"`
	Type          MessageFlowType `db:"type"`
	Slug          string          `db:"slug"`
	ValidateInput bool            `db:"validate_input"`
	Checkout      bool            `db:"checkout"`
	IsInput       bool
	IsReValidate  bool
}

type Message struct {
	ID           int64           `db:"id"`
	Slug         string          `db:"slug"`
	HeaderText   string          `db:"header_text"`
	PreviewURL   bool            `db:"preview_url"`
	BodyText     string          `db:"body_text"`
	FooterText   string          `db:"footer_text"`
	Button       string          `db:"button"`
	Type         MessageFlowType `db:"type"`
	WithMetadata bool            `db:"with_metadata"`
}

type MessageAction struct {
	Slug  string `db:"slug"`
	Title string `db:"title"`
	Desc  string `db:"description"`
}

type Profile struct {
	Nama         string
	TglLahir     string
	Hp           string
	Alamat       string
	JenisKelamin string
	Email        string
	Ktp          string
}

type Customer struct {
	ID             int64            `db:"id"`
	ClientID       int64            `db:"client_id"`
	WAID           string           `db:"wa_id"`
	Email          sql.Null[string] `db:"email"`
	FullName       sql.Null[string] `db:"full_name"`
	BirthDate      sql.Null[string] `db:"birth_date"`
	Address        sql.Null[string] `db:"address"`
	IdentityNumber sql.Null[string] `db:"identity_number"`
	IdentityType   sql.Null[string] `db:"identity_type"`
	Gender         sql.Null[string] `db:"gender"`
	Access         sql.Null[Access] `db:"access"`
	CreatedAt      string           `db:"created_at"`
	UpdatedAt      string           `db:"updated_at"`
}

type Payment struct {
	ID              int64               `db:"id"`
	CreatedAt       string              `db:"created_at"`
	UpdatedAt       string              `db:"updated_at"`
	ExpiredAt       sql.Null[string]    `db:"expired_at"`
	RefID           sql.Null[string]    `db:"ref_id"`
	ClientID        int64               `db:"client_id"`
	CustomerID      int64               `db:"customer_id"`
	PaymentType     sql.Null[string]    `db:"payment_type"`
	PaymentProvider sql.Null[string]    `db:"payment_provider"`
	PaymentCode     sql.Null[string]    `db:"payment_code"`
	PaymentRefID    sql.Null[string]    `db:"payment_ref_id"`
	PaymentItem     sql.Null[string]    `db:"payment_item"`
	Status          sql.Null[string]    `db:"status"`
	Amount          decimal.NullDecimal `db:"amount"`
	Fee             decimal.NullDecimal `db:"fee"`
}

type PaymentCustomer struct {
	ID         int64  `db:"id"`
	ClientID   int64  `db:"client_id"`
	WAPhoneID  string `db:"wa_phone_id"`
	WAID       string `db:"wa_id"`
	CustomerID int64  `db:"customer_id"`
	FullName   string `db:"full_name"`
}

type Session struct {
	ID        int64  `db:"id"`
	ClientID  int64  `db:"client_id"`
	WAID      string `db:"wa_id"`
	Access    Access `db:"access"`
	Seq       string `db:"seq"`
	Slug      string `db:"slug"`
	Input     string `db:"input"`
	CreatedAt string `db:"created_at"`
	ExpiredAt string `db:"expired_at"`
}

type MessageMetadata struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Invoice     string `json:"invoice"`
	ProductName string `json:"product_name"`
	CheckoutURL string `json:"checkout_url"`
	Amount      string `json:"amount"`
	ExpiryDate  string `json:"expiry_date"`
}

type Product struct {
	ID          int64           `db:"id"`
	Name        string          `db:"name"`
	Slug        string          `db:"slug"`
	Description sql.NullString  `db:"description"`
	Price       decimal.Decimal `db:"price"`
	Stock       decimal.Decimal `db:"stock"`
}

type (
	PaymentLinkCustomer struct {
		Email       string `json:"email"`
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		MobilePhone string `json:"mobilePhone"`
	}
	PaymentLinkOrder struct {
		ID          string `json:"id"`
		Amount      string `json:"amount"`
		Description string `json:"description"`
	}
	PaymentLinkCallbackURL struct {
		CallbackURL string `json:"callbackUrl"`
	}
	PaymentLink struct {
		Customer PaymentLinkCustomer    `json:"customer"`
		Order    PaymentLinkOrder       `json:"order"`
		URL      PaymentLinkCallbackURL `json:"url"`
	}
	PaymentLinkResponse struct {
		ResponseCode    string  `json:"responseCode"`
		ResponseMessage string  `json:"responseMessage"`
		ExpiryLink      string  `json:"expiryLink"`
		PaymentCode     string  `json:"paymentCode"`
		Appurl          string  `json:"appurl"`
		Imageurl        string  `json:"imageurl"`
		StringQr        string  `json:"stringQr"`
		Redirecturl     string  `json:"redirecturl"`
		ProcessingTime  float64 `json:"processingTime"`
		TraceID         string  `json:"traceId"`
	}
)

type FinpayCallback struct {
	Merchant struct {
		ID string `json:"id"`
	} `json:"merchant"`
	Customer struct {
		ID string `json:"id"`
	} `json:"customer"`
	Order struct {
		ID        string `json:"id"`
		Reference string `json:"reference"`
		Amount    int    `json:"amount"`
		Currency  string `json:"currency"`
	} `json:"order"`
	Transaction struct {
		AcquirerID any `json:"acquirerId"`
	} `json:"transaction"`
	SourceOfFunds struct {
		Type        string `json:"type"`
		PaymentCode string `json:"paymentCode"`
	} `json:"sourceOfFunds"`
	Meta struct {
		Data any `json:"data"`
	} `json:"meta"`
	Result struct {
		Payment struct {
			Status     string  `json:"status"`
			StatusDesc string  `json:"statusDesc"`
			UserDesc   string  `json:"userDesc"`
			Datetime   string  `json:"datetime"`
			Reference  string  `json:"reference"`
			Channel    string  `json:"channel"`
			Amount     float64 `json:"amount"`
		} `json:"payment"`
	} `json:"result"`
	Signature string `json:"signature"`
}
