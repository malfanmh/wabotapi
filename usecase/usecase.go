package usecase

import (
	"context"
	"github.com/malfanmh/wabotapi/model"
)

type Repository interface {
	GetClientByHash(ctx context.Context, hash string) (result model.Client, err error)
}

type WhatsAppRepository interface {
	SendTemplate(ctx context.Context, from, to string, template model.WATemplate, params map[string]interface{}) (string, error)
	SendText(ctx context.Context, from, to string, text string, params map[string]interface{}) (string, error)
}

type useCase struct {
	repo Repository
	wa   WhatsAppRepository
}

func New(repository Repository, wa WhatsAppRepository) *useCase {
	return &useCase{
		repo: repository,
		wa:   wa,
	}
}
