package myerrors

import (
	"fmt"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type ViolationError struct {
	Err  error
	User entity.User
}

func (ve *ViolationError) Error() string {
	return fmt.Sprintf("user with login %s already exists", ve.User.Login)
}

func (ve *ViolationError) Unwrap() error {
	return ve.Err
}

func NewViolationError(user entity.User, err error) error {

	return &ViolationError{
		User: user,
		Err:  err,
	}
}
