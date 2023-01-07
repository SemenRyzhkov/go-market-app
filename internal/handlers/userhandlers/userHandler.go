package userhandlers

import "net/http"

type UserHandler interface {
	Create(writer http.ResponseWriter, request *http.Request)
	Login(writer http.ResponseWriter, request *http.Request)
}
