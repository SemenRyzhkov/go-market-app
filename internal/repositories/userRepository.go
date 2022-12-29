package repositories

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type UserRepository interface {
	Save(ctx context.Context, login, password string) error
	FindByLogin(ctx context.Context, login string) (entity.User, error)
}
