package userservice

import (
	"context"
	"encoding/base64"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories"
)

var _ UserService = &userServiceImpl{}

type userServiceImpl struct {
	userRepository repositories.UserRepository
}

func (u userServiceImpl) Create(ctx context.Context, user entity.User) error {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(user.Password))

	return u.userRepository.Save(ctx, user.Login, encodedPassword)
}

func (u userServiceImpl) Login(ctx context.Context, user entity.User) error {
	foundUser, err := u.userRepository.FindByLogin(ctx, user.Login)
	if err != nil {
		return err
	}
	decodedPassword, err := base64.StdEncoding.DecodeString(foundUser.Password)
	if err != nil {
		return err
	}

	if user.Password != string(decodedPassword) {
		return &myerrors.InvalidPasswordError{Password: user.Password}
	}

	return nil
}

func New(userRepository repositories.UserRepository) UserService {
	return &userServiceImpl{
		userRepository,
	}
}
