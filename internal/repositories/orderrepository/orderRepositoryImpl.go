package orderrepository

import (
	"context"
	"database/sql"

	"github.com/omeid/pgerror"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
)

const (
	insertOrderQuery = "" +
		"INSERT INTO public.orders (number, status, uploaded_at , user_id) " +
		"VALUES ($1, $2, $3, $4)"
	findOrderByNumberQuery = "" +
		"SELECT number, user_id FROM public.orders " +
		"WHERE number=$1"
	getAllOrdersByUserIDQuery = "" +
		"SELECT number, status, " +
		"uploaded_at::timestamptz  " +
		"FROM public.orders " +
		"WHERE user_id=$1"
)

type orderRepositoryImpl struct {
	db *sql.DB
}

func New(db *sql.DB) OrderRepository {
	return &orderRepositoryImpl{
		db: db,
	}
}

func (r *orderRepositoryImpl) Save(ctx context.Context, order entity.Order) error {
	existedOrder, err := r.FindByNumber(ctx, order.Number)
	if orderWasSavedByOtherUser(err, existedOrder, order) {
		return myerrors.NewExistedOrderError(existedOrder.Number, existedOrder.UserID)
	}

	_, err = r.db.ExecContext(ctx, insertOrderQuery, order.Number, order.Status, order.UploadedAt, order.UserID)

	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return myerrors.NewOrderViolationError(order, err)
		}
	}
	return nil
}

func (r *orderRepositoryImpl) FindByNumber(ctx context.Context, order string) (entity.Order, error) {
	var existedOrder entity.Order
	row := r.db.QueryRowContext(ctx, findOrderByNumberQuery, order)
	err := row.Scan(&existedOrder.Number, &existedOrder.UserID)
	if err != nil && err != sql.ErrNoRows {
		return entity.Order{}, err
	}
	return existedOrder, nil
}

func (r *orderRepositoryImpl) GetAllByUserID(ctx context.Context, userID string) ([]entity.Order, error) {
	orderList := make([]entity.Order, 0)

	rows, err := r.db.QueryContext(ctx, getAllOrdersByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o entity.Order
		err = rows.Scan(&o.Number, &o.Status, &o.UploadedAt)
		if err != nil {
			return nil, err
		}

		orderList = append(orderList, o)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return orderList, nil
}

func orderWasSavedByOtherUser(err error, existedOrder entity.Order, order entity.Order) bool {
	return err == nil && len(existedOrder.Number) > 0 && existedOrder.UserID != order.UserID
}
