package myerrors

import (
	"fmt"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

type OrderViolationError struct {
	Err   error
	Order entity.Order
}

func (oe *OrderViolationError) Error() string {
	return fmt.Sprintf("order with number %s already exists", oe.Order.Number)
}

func (oe *OrderViolationError) Unwrap() error {
	return oe.Err
}

func NewOrderViolationError(order entity.Order, err error) error {

	return &OrderViolationError{
		Order: order,
		Err:   err,
	}
}
