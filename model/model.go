package model

type Client struct {
	ID         int64  `db:"id"`
	MerchantID int64  `db:"merchant_id"`
	Name       string `db:"name"`
	Hash       string `db:"hash"`
	Token      string `db:"token"`
}
