package withdrawservice

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type WithdrawService interface {
	Create(ctx context.Context, withdrawRequest entity.WithdrawRequest, userID string) error
	GetUserBalance(ctx context.Context, userID string) (entity.BalanceRequest, error)
	GetAllByUserID(ctx context.Context, userID string) ([]entity.WithdrawDTO, error)
}
