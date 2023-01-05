package service

import (
	"strconv"

	"github.com/theplant/luhn"

	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
)

func CheckNumberByLuhnAlgorithm(order string) (int, error) {
	orderNumber, err := strconv.Atoi(order)
	if err != nil {
		return 0, err
	}
	if !luhn.Valid(orderNumber) {
		return 0, myerrors.NewInvalidOrderNumberFormatError(orderNumber)
	}
	return orderNumber, nil
}
