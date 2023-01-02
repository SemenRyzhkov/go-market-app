package orderservice

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type OrderService interface {
	Create(ctx context.Context, order, userID string) error
	GetAllByUserID(ctx context.Context, userID string) ([]entity.OrderDTO, error)
}
