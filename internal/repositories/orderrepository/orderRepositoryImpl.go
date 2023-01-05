package orderrepository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/omeid/pgerror"

	"github.com/SemenRyzhkov/go-market-app/internal/client"
	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/service/orderservice/ordermapper"
)

const (
	tickerDuration   = 5 * time.Second
	updateOrderQuery = "" +
		"UPDATE public.orders " +
		"SET status = $1, " +
		"    accrual = $2 " +
		"WHERE number = $3"
	insertOrderQuery = "" +
		"INSERT INTO public.orders (number, status, accrual, uploaded_at, user_id) " +
		"VALUES ($1, $2, $3, $4, $5)"
	findOrderByNumberQuery = "" +
		"SELECT number, user_id FROM public.orders " +
		"WHERE number=$1"
	getAllOrdersByUserIDQuery = "" +
		"SELECT number, status, accrual, " +
		"uploaded_at::timestamptz " +
		"FROM public.orders " +
		"WHERE user_id=$1"
	getTotalAccrualByUserIDQuery = "" +
		"SELECT SUM(accrual) " +
		"FROM public.orders " +
		"WHERE user_id = $1"
	getOrderNumbersWithStatusNewOrProcessingQuery = "" +
		"SELECT number " +
		"FROM public.orders " +
		"WHERE status IN (1, 2)"
)

var (
	_                      OrderRepository = &orderRepositoryImpl{}
	ErrRepositoryIsClosing                 = errors.New("repository is closing")
	ErrSchedulerIsClosing                  = errors.New("scheduler is closing")
)

type orderRepositoryImpl struct {
	db                 *sql.DB
	ticker             *time.Ticker
	done               chan struct{}
	updatingOrderQueue chan int
	wg                 sync.WaitGroup
	once               sync.Once
	client             *client.Client
}

func New(db *sql.DB, host string, duration time.Duration) (OrderRepository, error) {
	repo := orderRepositoryImpl{
		db:                 db,
		ticker:             time.NewTicker(tickerDuration),
		done:               make(chan struct{}),
		updatingOrderQueue: make(chan int),
		client:             client.NewClient(host, duration),
	}
	err := repo.runSchedulerForLoadNotUpdatedOrders()
	if err != nil {
		return nil, err
	}
	repo.runUpdatingStatusWorkerPool()
	return &repo, nil
}

func (r *orderRepositoryImpl) runSchedulerForLoadNotUpdatedOrders() error {
	var err error
	go func() {
		for {
			select {
			case <-r.done:
				err = ErrSchedulerIsClosing
			case <-r.ticker.C:
				getErr := r.getOrderNumbersWithStatusNewOrProcessingAndAddItToUpdatingQueue(context.Background())
				if err != nil {
					err = getErr
					return
				}
			}
		}
	}()
	return err
}

func (r *orderRepositoryImpl) runUpdatingStatusWorkerPool() {
	for i := 0; i < 10; i++ {
		r.wg.Add(1)
		go func() {
			defer r.wg.Done()
			for {
				select {
				case <-r.done:
					log.Println("Exiting")
					return
				case num, ok := <-r.updatingOrderQueue:
					if !ok {
						return
					}
					orderResponse, err := r.client.GetOrderResponse(num)
					if err != nil {
						log.Printf("get order response error %v", err)
						return
					}
					order, err := ordermapper.MapOrderResponseToOrder(orderResponse)
					if err != nil {
						log.Printf("mapper error %v", err)
						return
					}
					_, err = r.db.ExecContext(context.Background(), updateOrderQuery, order.Status, order.Accrual, order.Number)
					if err != nil {
						log.Printf("update order error %v", err)
						return
					}
				}
			}
		}()
	}
}

func (r *orderRepositoryImpl) StopSchedulerAndWorkerPool() {
	r.once.Do(func() {
		close(r.done)
		close(r.updatingOrderQueue)
	})
	r.wg.Wait()
	r.ticker.Stop()
}

func (r *orderRepositoryImpl) Save(ctx context.Context, order entity.Order) error {
	existedOrder, err := r.FindByNumber(ctx, order.Number)
	if orderWasSavedByOtherUser(err, existedOrder, order) {
		return myerrors.NewExistedOrderError(existedOrder.Number, existedOrder.UserID)
	}

	_, err = r.db.ExecContext(ctx, insertOrderQuery, order.Number, order.Status, 0, order.UploadedAt, order.UserID)

	if err != nil {
		if e := pgerror.UniqueViolation(err); e != nil {
			return myerrors.NewOrderViolationError(order, err)
		}
	}
	return nil
}

func (r *orderRepositoryImpl) FindByNumber(ctx context.Context, number int) (entity.Order, error) {
	var existedOrder entity.Order
	row := r.db.QueryRowContext(ctx, findOrderByNumberQuery, number)
	err := row.Scan(&existedOrder.Number, &existedOrder.UserID)
	if err != nil && err != sql.ErrNoRows {
		return entity.Order{}, err
	}
	return existedOrder, nil
}

func (r *orderRepositoryImpl) GetAllByUserID(ctx context.Context, userID string) ([]entity.Order, error) {
	orderList := make([]entity.Order, 0)

	rows, err := r.db.QueryContext(ctx, getAllOrdersByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var o entity.Order
		err = rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}

		orderList = append(orderList, o)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	sort.Slice(orderList, func(i, j int) bool { return orderList[j].UploadedAt.Before(orderList[i].UploadedAt) })

	return orderList, nil
}

func (r *orderRepositoryImpl) GetTotalAccrualByUserID(ctx context.Context, userID string) (float64, error) {
	var totalAccrual float64
	row := r.db.QueryRowContext(ctx, getTotalAccrualByUserIDQuery, userID)
	err := row.Scan(&totalAccrual)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return totalAccrual, nil
}

func (r *orderRepositoryImpl) getOrderNumbersWithStatusNewOrProcessingAndAddItToUpdatingQueue(ctx context.Context) error {
	rows, err := r.db.QueryContext(ctx, getOrderNumbersWithStatusNewOrProcessingQuery)
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var num int
		err = rows.Scan(&num)
		if err != nil {
			return err
		}

		err := r.addNumberToUpdatingQueue(num)
		if err != nil {
			return err
		}
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepositoryImpl) addNumberToUpdatingQueue(number int) error {
	select {
	case <-r.done:
		return ErrRepositoryIsClosing
	case r.updatingOrderQueue <- number:
		return nil
	}
}

func orderWasSavedByOtherUser(err error, existedOrder entity.Order, order entity.Order) bool {
	return err == nil && existedOrder.Number > 0 && existedOrder.UserID != order.UserID
}
