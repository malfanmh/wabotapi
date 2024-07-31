package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/malfanmh/wabotapi/model"
)

type mysqlRepository struct {
	db *sqlx.DB
}

func NewMysql(db *sqlx.DB) *mysqlRepository {
	return &mysqlRepository{db}
}

func (r *mysqlRepository) GetClientByHash(ctx context.Context, hash string) (result model.Client, err error) {
	q := `SELECT id, name, hash, token FROM clients WHERE hash = ?`
	err = r.db.GetContext(ctx, &result, r.db.Rebind(q), hash)
	return
}
