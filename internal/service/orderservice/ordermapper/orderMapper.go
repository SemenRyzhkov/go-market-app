package ordermapper

import (
	"time"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

func MapToOrder(order, userID string) entity.Order {
	return entity.Order{
		Number:     order,
		UploadedAt: time.Now(),
		Status:     entity.NEW,
		UserID:     userID,
	}
}

func MapOrderListToOrderDTOList(orderList []entity.Order) []entity.OrderDTO {
	orderDTOList := make([]entity.OrderDTO, 0)
	for _, o := range orderList {

		dto := entity.OrderDTO{
			Number:     o.Number,
			Status:     o.Status.String(),
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		}
		orderDTOList = append(orderDTOList, dto)
	}
	return orderDTOList
}
