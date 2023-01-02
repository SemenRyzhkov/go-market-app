package userservice

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type UserService interface {
	Create(ctx context.Context, user entity.UserRequest, userID string) error
	Login(ctx context.Context, user entity.UserRequest) (string, error)
}
