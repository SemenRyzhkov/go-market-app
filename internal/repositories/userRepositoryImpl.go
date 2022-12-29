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
	_                      UserRepository = &userRepositoryImpl{}
	ErrRepositoryIsClosing                = errors.New("repository is closing")
)

const (
	initDBQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.users (" +
		"id serial primary key, " +
		"login varchar(45) unique not null, " +
		"password varchar(45) not null, " +
		")"
	createUserLoginIndex = "" +
		"CREATE INDEX IF NOT EXISTS user_login_index " +
		"ON public.users (login)"
	insertUserQuery = "" +
		"INSERT INTO public.users (login, password) " +
		"VALUES ($1, $2)"
	findUserByLoginQuery = "" +
		"SELECT login, password FROM public.users " +
		"WHERE login=$1"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func New(dbAddress string) (UserRepository, error) {
	db, err := initDB(dbAddress)
	if err != nil {
		return nil, err
	}
	return &userRepositoryImpl{
		db: db,
	}, nil

}

func (u *userRepositoryImpl) Save(ctx context.Context, login, password string) error {
	_, err := u.db.ExecContext(ctx, insertUserQuery, login, password)
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return myerrors.NewViolationError(entity.User{Login: login, Password: password}, err)
		}
	}
	return nil
}

func (u *userRepositoryImpl) FindByLogin(ctx context.Context, login string) (entity.User, error) {
	var user entity.User
	row := u.db.QueryRowContext(ctx, findUserByLoginQuery, login)
	err := row.Scan(&user.Login, &user.Password)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
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
	_, createTableErr := db.Exec(initDBQuery)
	if createTableErr != nil {
		return createTableErr
	}
	_, createIndexErr := db.Exec(createUserLoginIndex)
	if createIndexErr != nil {
		return createIndexErr
	}
	return nil
}
