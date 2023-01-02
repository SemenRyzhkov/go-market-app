package repositories

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type Repository interface {
	SaveUser(ctx context.Context, userID, login, password string) error
	FindUserByLogin(ctx context.Context, login string) (entity.UserDTO, error)
	FindOrderByNumber(ctx context.Context, order string) (entity.Order, error)
	SaveOrder(ctx context.Context, order entity.Order) error
	GetAllOrdersByUserID(ctx context.Context, userID string) ([]entity.Order, error)
}
