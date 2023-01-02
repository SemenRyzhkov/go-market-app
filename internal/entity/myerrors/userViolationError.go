package myerrors

import (
	"fmt"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type UserViolationError struct {
	Err  error
	User entity.UserRequest
}

func (ve *UserViolationError) Error() string {
	return fmt.Sprintf("user with login %s already exists", ve.User.Login)
}

func (ve *UserViolationError) Unwrap() error {
	return ve.Err
}

func NewUserViolationError(user entity.UserRequest, err error) error {

	return &UserViolationError{
		User: user,
		Err:  err,
	}
}
