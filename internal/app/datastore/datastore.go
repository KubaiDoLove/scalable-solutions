package datastore

import (
	"context"
	"github.com/KubaiDoLove/scalable-solutions/internal/app/models"
	"github.com/google/uuid"
)

type DataStore interface {
	CreateOrder(ctx context.Context, order *models.Order) error
	DisableOrder(ctx context.Context, id uuid.UUID) error
	OrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	MatchOrder(ctx context.Context, order *models.Order) ([]models.Order, error)
	MarketDataSnapshot(ctx context.Context) (*models.MarketDataSnapshot, error)
}
