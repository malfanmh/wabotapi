package usecase

import (
	"context"
	"github.com/malfanmh/wabotapi/model"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type Repository interface {
	GetClientByHash(ctx context.Context, hash string) (result model.Client, err error)
}

type WhatsAppRepository interface {
	Send(ctx context.Context, from, to string, msgType, jsonBody string) (string, error)
	SendTemplate(ctx context.Context, from, to string, template model.WATemplate, params map[string]interface{}) (string, error)
	SendText(ctx context.Context, from, to string, text string, params map[string]interface{}) (string, error)
}

type useCase struct {
	repo    Repository
	wa      WhatsAppRepository
	members cmap.ConcurrentMap[string, model.Member]
}

func New(repository Repository, wa WhatsAppRepository) *useCase {
	return &useCase{
		repo:    repository,
		wa:      wa,
		members: cmap.New[model.Member](),
	}
}
