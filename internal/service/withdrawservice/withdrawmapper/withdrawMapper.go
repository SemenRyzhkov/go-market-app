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
		Current:   totalUserAccrual,
		Withdrawn: totalUserWithdraw,
	}
}

func MapOrderListToOrderDTOList(orderList []entity.Order) []entity.OrderDTO {
	orderDTOList := make([]entity.OrderDTO, 0)
	for _, o := range orderList {

		dto := entity.OrderDTO{
			Number:     strconv.Itoa(o.Number),
			Status:     o.Status.String(),
			Accrual:    o.Accrual,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		}
		orderDTOList = append(orderDTOList, dto)
	}
	return orderDTOList
}
