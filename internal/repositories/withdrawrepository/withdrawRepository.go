package withdrawrepository

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type WithdrawRepository interface {
	Save(ctx context.Context, withdraw entity.Withdraw) error
	GetTotalWithdrawByUserID(ctx context.Context, userID string) (float32, error)
	GetAllByUserID(ctx context.Context, userID string) ([]entity.Withdraw, error)
}
