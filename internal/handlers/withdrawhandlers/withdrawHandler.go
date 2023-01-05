package withdrawhandlers

import "net/http"

type WithdrawHandler interface {
	Create(writer http.ResponseWriter, request *http.Request)
	GetAll(writer http.ResponseWriter, request *http.Request)
}
