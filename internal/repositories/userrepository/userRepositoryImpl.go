package userrepository

import (
	"context"
	"database/sql"

	"github.com/omeid/pgerror"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
)

const (
	insertUserQuery = "" +
		"INSERT INTO public.users (id, login, password) " +
		"VALUES ($1, $2, $3)"
	findUserByLoginQuery = "" +
		"SELECT id, login, password FROM public.users " +
		"WHERE login=$1"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func New(db *sql.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Save(ctx context.Context, userID, login, password string) error {
	_, err := r.db.ExecContext(ctx, insertUserQuery, userID, login, password)
	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return myerrors.NewUserViolationError(entity.UserRequest{Login: login, Password: password}, err)
		}
	}
	return nil
}

func (r *userRepositoryImpl) FindByLogin(ctx context.Context, login string) (entity.UserDTO, error) {
	var user entity.UserDTO
	row := r.db.QueryRowContext(ctx, findUserByLoginQuery, login)
	err := row.Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return entity.UserDTO{}, err
	}
	return user, nil
}
