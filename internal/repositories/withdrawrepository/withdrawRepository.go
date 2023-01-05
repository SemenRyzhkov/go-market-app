package withdrawrepository

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type WithdrawRepository interface {
	Save(ctx context.Context, withdraw entity.Withdraw) error
}
