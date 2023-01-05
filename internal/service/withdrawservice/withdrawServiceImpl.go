package withdrawservice

import (
	"context"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories/orderrepository"
	"github.com/SemenRyzhkov/go-market-app/internal/repositories/withdrawrepository"
	"github.com/SemenRyzhkov/go-market-app/internal/service"
	"github.com/SemenRyzhkov/go-market-app/internal/service/withdrawservice/withdrawmapper"
)

var _ WithdrawService = &withdrawServiceImpl{}

type withdrawServiceImpl struct {
	withdrawRepository withdrawrepository.WithdrawRepository
	orderRepository    orderrepository.OrderRepository
}

func (w *withdrawServiceImpl) Create(ctx context.Context, withdrawRequest entity.WithdrawRequest, userID string) error {
	order, err := service.CheckNumberByLuhnAlgorithm(withdrawRequest.Order)
	if err != nil {
		return err
	}
	totalUserAccrual, err := w.orderRepository.GetTotalAccrualByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if totalUserAccrual < withdrawRequest.Sum {
		return myerrors.NewLimitExceededError(withdrawRequest.Sum, totalUserAccrual, userID)
	}

	return w.withdrawRepository.Save(ctx, withdrawmapper.MapToWithdraw(withdrawRequest, order, userID))
}

func (w *withdrawServiceImpl) GetAllByUserID(ctx context.Context, userID string) ([]entity.OrderDTO, error) {

	return []entity.OrderDTO{}, nil
}

func New(withdrawRepository withdrawrepository.WithdrawRepository, orderRepository orderrepository.OrderRepository) WithdrawService {
	return &withdrawServiceImpl{
		withdrawRepository: withdrawRepository,
		orderRepository:    orderRepository,
	}
}
