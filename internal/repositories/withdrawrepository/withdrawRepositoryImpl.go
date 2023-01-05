package withdrawrepository

import (
	"context"
	"database/sql"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

const (
	insertWithdrawQuery = "" +
		"INSERT INTO public.withdraw (number, sum, processed_at, user_id) " +
		"VALUES ($1, $2, $3, $4)"
	getTotalWithdrawByUserIDQuery = "" +
		"SELECT SUM(sum) " +
		"FROM public.withdraw " +
		"WHERE user_id = $1"

	findOrderByNumberQuery = "" +
		"SELECT number, user_id FROM public.orders " +
		"WHERE number=$1"
	getAllOrdersByUserIDQuery = "" +
		"SELECT number, status, accrual, " +
		"uploaded_at::timestamptz " +
		"FROM public.orders " +
		"WHERE user_id=$1"
	getOrderNumbersWithStatusNewOrProcessingQuery = "" +
		"SELECT number " +
		"FROM public.orders " +
		"WHERE status IN (1, 2)"
)

var (
	_ WithdrawRepository = &withdrawRepositoryImpl{}
)

type withdrawRepositoryImpl struct {
	db *sql.DB
}

func New(db *sql.DB) WithdrawRepository {
	return &withdrawRepositoryImpl{
		db: db,
	}

}

func (w *withdrawRepositoryImpl) Save(ctx context.Context, withdraw entity.Withdraw) error {
	_, err := w.db.ExecContext(ctx, insertWithdrawQuery, withdraw.Order, withdraw.Sum, withdraw.ProcessedAt, withdraw.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (w *withdrawRepositoryImpl) GetTotalWithdrawByUserID(ctx context.Context, userID string) (float32, error) {
	var totalAccrual float32
	row := w.db.QueryRowContext(ctx, getTotalWithdrawByUserIDQuery, userID)
	err := row.Scan(&totalAccrual)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return totalAccrual, nil
}
