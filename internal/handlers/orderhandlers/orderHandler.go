package orderhandlers

import "net/http"

type OrderHandler interface {
	Create(writer http.ResponseWriter, request *http.Request)
	GetAll(writer http.ResponseWriter, request *http.Request)
}
