package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/go-market-app/internal/handlers/orderhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/handlers/userhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/router/middleware"
)

const (
	createUserPath  = "/api/user/register"
	createOrderPath = "/api/user/orders"
	loginUserPath   = "/api/user/login"
)

func NewRouter(h userhandlers.UserHandler, o orderhandlers.OrderHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.DecompressRequest, middleware.CompressResponse, middleware.LoggingMiddleware)
	r.Post(createUserPath, h.Create)
	r.Post(createOrderPath, o.Create)
	r.Post(loginUserPath, h.Login)
	r.Get(createOrderPath, o.GetAll)

	return r
}
