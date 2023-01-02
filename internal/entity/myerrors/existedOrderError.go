package myerrors

import (
	"fmt"
)

type ExistedOrderError struct {
	Number string
	UserID string
}

func (eo *ExistedOrderError) Error() string {
	return fmt.Sprintf("order with number %s was saved by user with ID %s", eo.Number, eo.UserID)
}

func NewExistedOrderError(number, userID string) error {
	return &ExistedOrderError{
		Number: number,
		UserID: userID,
	}
}
