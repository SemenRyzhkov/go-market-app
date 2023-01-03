package orderrepository

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type OrderRepository interface {
	FindByNumber(ctx context.Context, order string) (entity.Order, error)
	Save(ctx context.Context, order entity.Order) error
	GetAllByUserID(ctx context.Context, userID string) ([]entity.Order, error)
}
