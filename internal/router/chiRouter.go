package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/SemenRyzhkov/go-market-app/internal/handlers/userhandlers"
	"github.com/SemenRyzhkov/go-market-app/internal/router/middleware"
)

const (
	createUserPath = "/api/user/register"
	loginUserPath  = "/api/user/login"
)

func NewRouter(h userhandlers.UserHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.DecompressRequest, middleware.CompressResponse, middleware.LoggingMiddleware)
	r.Post(createUserPath, h.Create)
	r.Post(loginUserPath, h.Login)

	return r
}
