package repositories

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/omeid/pgerror"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
)

var (
	_                      Repository = &repositoryImpl{}
	ErrRepositoryIsClosing            = errors.New("repository is closing")
)

const (
	initUsersTableQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.users (" +
		"id varchar(45) primary key, " +
		"login varchar(45) unique not null, " +
		"password varchar(45) not null" +
		")"
	initOrdersTableQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.orders (" +
		"number varchar(45) primary key, " +
		"status int2 not null, " +
		"accrual decimal, " +
		"uploaded_at timestamptz not null, " +
		"user_id varchar(45) references public.users (id)" +
		")"
	setTimeZoneQuery = "" +
		"set timezone = 'Europe/Moscow'"
	createUserLoginIndex = "" +
		"CREATE INDEX IF NOT EXISTS user_login_index " +
		"ON public.users (login)"
	insertUserQuery = "" +
		"INSERT INTO public.users (id, login, password) " +
		"VALUES ($1, $2, $3)"
	insertOrderQuery = "" +
		"INSERT INTO public.orders (number, status, uploaded_at , user_id) " +
		"VALUES ($1, $2, $3, $4)"
	findUserByLoginQuery = "" +
		"SELECT id, login, password FROM public.users " +
		"WHERE login=$1"
	findOrderByNumberQuery = "" +
		"SELECT number, user_id FROM public.orders " +
		"WHERE number=$1"
	getAllOrdersByUserIDQuery = "" +
		"SELECT number, status, " +
		"uploaded_at::timestamptz  " +
		"FROM public.orders " +
		"WHERE user_id=$1"
)

type repositoryImpl struct {
	db *sql.DB
}

func New(dbAddress string) (Repository, error) {
	db, err := initDB(dbAddress)
	if err != nil {
		return nil, err
	}
	return &repositoryImpl{
		db: db,
	}, nil

}

func (r *repositoryImpl) SaveUser(ctx context.Context, userID, login, password string) error {
	_, err := r.db.ExecContext(ctx, insertUserQuery, userID, login, password)
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return myerrors.NewUserViolationError(entity.UserRequest{Login: login, Password: password}, err)
		}
	}
	return nil
}

func (r *repositoryImpl) SaveOrder(ctx context.Context, order entity.Order) error {
	existedOrder, err := r.FindOrderByNumber(ctx, order.Number)
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

func (r *repositoryImpl) FindOrderByNumber(ctx context.Context, order string) (entity.Order, error) {
	var existedOrder entity.Order
	row := r.db.QueryRowContext(ctx, findOrderByNumberQuery, order)
	err := row.Scan(&existedOrder.Number, &existedOrder.UserID)
	if err != nil && err != sql.ErrNoRows {
		return entity.Order{}, err
	}
	return existedOrder, nil
}

func orderWasSavedByOtherUser(err error, existedOrder entity.Order, order entity.Order) bool {
	return err == nil && len(existedOrder.Number) > 0 && existedOrder.UserID != order.UserID
}

func (r *repositoryImpl) FindUserByLogin(ctx context.Context, login string) (entity.UserDTO, error) {
	var user entity.UserDTO
	row := r.db.QueryRowContext(ctx, findUserByLoginQuery, login)
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return entity.UserDTO{}, err
	}
	return user, nil
}

func (r *repositoryImpl) GetAllOrdersByUserID(ctx context.Context, userID string) ([]entity.Order, error) {
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

func initDB(dbAddress string) (*sql.DB, error) {
	db, connectionErr := sql.Open("postgres", dbAddress)
	if connectionErr != nil {
		log.Println(connectionErr)
		return nil, connectionErr
	}
	createTableErr := createTableIfNotExists(db)
	if createTableErr != nil {
		log.Println(createTableErr)
		return nil, createTableErr
	}
	return db, nil
}

func createTableIfNotExists(db *sql.DB) error {
	_, createUserTableErr := db.Exec(initUsersTableQuery)
	if createUserTableErr != nil {
		return createUserTableErr
	}
	_, createOrderTableErr := db.Exec(initOrdersTableQuery)
	if createOrderTableErr != nil {
		return createUserTableErr
	}
	_, createIndexErr := db.Exec(createUserLoginIndex)
	if createIndexErr != nil {
		return createIndexErr
	}
	_, setTimeZoneErr := db.Exec(setTimeZoneQuery)
	if setTimeZoneErr != nil {
		return setTimeZoneErr
	}
	return nil
}
