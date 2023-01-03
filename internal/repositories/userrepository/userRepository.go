package userrepository

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type UserRepository interface {
	Save(ctx context.Context, userID, login, password string) error
	FindByLogin(ctx context.Context, login string) (entity.UserDTO, error)
}
