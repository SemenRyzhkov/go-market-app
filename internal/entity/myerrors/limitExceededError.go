package myerrors

import (
	"fmt"
)

type LimitExceededError struct {
	Sum          float64
	TotalAccrual float64
	UserID       string
}

func (le *LimitExceededError) Error() string {
	return fmt.Sprintf("sum %f exceededs limit, for user %s accessible limit is %f", le.Sum, le.UserID, le.TotalAccrual)
}

func NewLimitExceededError(sum float64, totalAccrual float64, userID string) error {
	return &LimitExceededError{
		Sum:          sum,
		TotalAccrual: totalAccrual,
		UserID:       userID,
	}
}
