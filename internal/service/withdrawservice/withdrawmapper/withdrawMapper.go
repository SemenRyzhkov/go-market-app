package withdrawmapper

import (
	"strconv"
	"time"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

func MapToWithdraw(withdrawRequest entity.WithdrawRequest, orderNumber int, userID string) entity.Withdraw {
	return entity.Withdraw{
		Order:       orderNumber,
		Sum:         withdrawRequest.Sum,
		ProcessedAt: time.Now(),
		UserID:      userID,
	}
}

func MapToBalanceRequest(totalUserAccrual float32, totalUserWithdraw float32) entity.BalanceRequest {
	return entity.BalanceRequest{
		Current:   totalUserAccrual - totalUserWithdraw,
		Withdrawn: totalUserWithdraw,
	}
}

func MapWithdrawListToWithdrawDTOList(withdrawList []entity.Withdraw) []entity.WithdrawDTO {
	withdrawDTOList := make([]entity.WithdrawDTO, 0)
	for _, w := range withdrawList {

		dto := entity.WithdrawDTO{
			Order:       strconv.Itoa(w.Order),
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
		}
		withdrawDTOList = append(withdrawDTOList, dto)
	}
	return withdrawDTOList
}
