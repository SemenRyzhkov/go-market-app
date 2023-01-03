package ordermapper

import (
	"strconv"
	"time"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

func MapToOrder(number int, userID string) entity.Order {
	return entity.Order{
		Number:     number,
		UploadedAt: time.Now(),
		Status:     entity.NEW,
		UserID:     userID,
	}
}

func MapOrderListToOrderDTOList(orderList []entity.Order) []entity.OrderDTO {
	orderDTOList := make([]entity.OrderDTO, 0)
	for _, o := range orderList {

		dto := entity.OrderDTO{
			Number:     string(o.Number),
			Status:     o.Status.String(),
			Accrual:    o.Accrual,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		}
		orderDTOList = append(orderDTOList, dto)
	}
	return orderDTOList
}

func MapOrderResponseToOrder(orderResponse entity.OrderResponse) (entity.Order, error) {
	number, err := strconv.Atoi(orderResponse.Order)
	if err != nil {
		return entity.Order{}, err

	}
	order := entity.Order{
		Number: number,
	}
	if orderResponse.Status == "REGISTERED" || orderResponse.Status == "PROCESSING" {
		order.Status = entity.OrderStatus(2)
	} else if orderResponse.Status == "INVALID" {
		order.Status = entity.OrderStatus(3)
	} else {
		order.Status = entity.OrderStatus(4)
		order.Accrual = orderResponse.Accrual
	}
	return order, nil
}
