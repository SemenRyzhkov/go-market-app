package userservice

import (
	"context"
	"encoding/base64"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories/userrepository"
)

var _ UserService = &userServiceImpl{}

type userServiceImpl struct {
	userRepository userrepository.UserRepository
}

func (u userServiceImpl) Create(ctx context.Context, user entity.UserRequest, userID string) error {
	encodedPassword := base64.StdEncoding.EncodeToString([]byte(user.Password))

	return u.userRepository.Save(ctx, userID, user.Login, encodedPassword)
}

func (u userServiceImpl) Login(ctx context.Context, user entity.UserRequest) (string, error) {
	foundUser, err := u.userRepository.FindByLogin(ctx, user.Login)
	if err != nil {
		return "", err
	}
	decodedPassword, err := base64.StdEncoding.DecodeString(foundUser.Password)
	if err != nil {
		return "", err
	}

	if user.Password != string(decodedPassword) {
		return "", &myerrors.InvalidPasswordError{Password: user.Password}
	}

	return foundUser.ID, nil
}

func New(userRepository userrepository.UserRepository) UserService {
	return &userServiceImpl{
		userRepository,
	}
}
