package model

import "database/sql"

type Access int

const (
	AccessPublic Access = iota
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
	ID         int64  `db:"id"`
	MerchantID int64  `db:"merchant_id"`
	Name       string `db:"name"`
	Hash       string `db:"hash"`
	Token      string `db:"token"`
	WAPhone    string `db:"wa_phone"`
	WAPhoneID  string `db:"wa_phone_id"`
}

type MessageFlowType string

const (
	MessageFlowText   MessageFlowType = "text"
	MessageFlowList   MessageFlowType = "list"
	MessageFlowButton MessageFlowType = "button"
)

var mapType = map[MessageFlowType]WAMessageType{
	MessageFlowButton: WAMessageTypeInteractive,
	MessageFlowList:   WAMessageTypeInteractive,
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
	IsInput       bool
	IsReValidate  bool
}

type Message struct {
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
	Status         sql.Null[Access] `db:"status"`
	CreatedAt      string           `db:"created_at"`
	UpdatedAt      string           `db:"updated_at"`
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
	Name string `json:"name"`
}
