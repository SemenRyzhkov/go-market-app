package withdrawrepository

import (
	"context"
	"database/sql"
	"sort"

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
	getAllWithdrawsByUserIDQuery = "" +
		"SELECT number, sum, " +
		"processed_at::timestamptz " +
		"FROM public.withdraw " +
		"WHERE user_id=$1"
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
	var totalWithdraw sql.NullFloat64
	row := w.db.QueryRowContext(ctx, getTotalWithdrawByUserIDQuery, userID)
	err := row.Scan(&totalWithdraw)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if totalWithdraw.Valid {
		return float32(totalWithdraw.Float64), nil
	} else {
		return 0.0, nil
	}
}

func (w *withdrawRepositoryImpl) GetAllByUserID(ctx context.Context, userID string) ([]entity.Withdraw, error) {
	withdrawList := make([]entity.Withdraw, 0)

	rows, err := w.db.QueryContext(ctx, getAllWithdrawsByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var w entity.Withdraw
		err = rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}

		withdrawList = append(withdrawList, w)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	sort.Slice(withdrawList, func(i, j int) bool { return withdrawList[j].ProcessedAt.Before(withdrawList[i].ProcessedAt) })

	return withdrawList, nil
}
