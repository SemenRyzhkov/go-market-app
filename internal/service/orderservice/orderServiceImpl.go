package orderservice

import (
	"context"
	"strconv"

	"github.com/theplant/luhn"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories/orderrepository"
	"github.com/SemenRyzhkov/go-market-app/internal/service/orderservice/ordermapper"
)

var _ OrderService = &orderServiceImpl{}

type orderServiceImpl struct {
	orderRepository orderrepository.OrderRepository
}

func (o orderServiceImpl) Create(ctx context.Context, order string, userID string) error {
	err := checkNumberByLuhnAlgorithm(order)
	if err != nil {
		return err
	}
	return o.orderRepository.Save(ctx, ordermapper.MapToOrder(order, userID))
}

func (o orderServiceImpl) GetAllByUserID(ctx context.Context, userID string) ([]entity.OrderDTO, error) {
	orderList, err := o.orderRepository.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return ordermapper.MapOrderListToOrderDTOList(orderList), nil
}

func New(orderRepository orderrepository.OrderRepository) OrderService {
	return &orderServiceImpl{
		orderRepository,
	}
}

func checkNumberByLuhnAlgorithm(order string) error {
	orderNumber, err := strconv.Atoi(order)
	if err != nil {
		return err
	}
	if !luhn.Valid(orderNumber) {
		return &myerrors.InvalidOrderNumberFormatError{Order: order}
	}
	return nil
}