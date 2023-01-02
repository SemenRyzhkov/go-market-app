package myerrors

import (
	"fmt"
)

type InvalidOrderNumberFormatError struct {
	Order string
}

func (iof *InvalidOrderNumberFormatError) Error() string {
	return fmt.Sprintf("order with number %s has invalid type", iof.Order)
}

func NewInvalidOrderNumberFormatError(order string) error {
	return &InvalidOrderNumberFormatError{
		Order: order,
	}
}
