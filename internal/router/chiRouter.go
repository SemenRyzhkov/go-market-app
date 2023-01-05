package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/go-market-app/internal/handlers/orderhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/handlers/userhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/handlers/withdrawhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/router/middleware"
)

const (
	createUserPath          = "/api/user/register"
	createOrderPath         = "/api/user/orders"
	createWithdrawPath      = "/api/user/balance/withdraw"
	loginUserPath           = "/api/user/login"
	getUserBalancePath      = "/api/user/balance"
	getAllUserWithdrawsPath = "/api/user/withdrawals"
)

func NewRouter(h userhandlers.UserHandler, o orderhandlers.OrderHandler, w withdrawhandlers.WithdrawHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.DecompressRequest, middleware.CompressResponse, middleware.LoggingMiddleware)
	r.Post(createUserPath, h.Create)
	r.Post(createOrderPath, o.Create)
	r.Post(createWithdrawPath, w.Create)
	r.Post(loginUserPath, h.Login)
	r.Get(createOrderPath, o.GetAll)
	r.Get(getUserBalancePath, w.GetUserBalance)
	r.Get(getAllUserWithdrawsPath, w.GetAll)

	return r
}
