package withdrawhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/go-market-app/internal/service/withdrawservice"
)

type withdrawHandlerImpl struct {
	withdrawService withdrawservice.WithdrawService
	cookieService   cookieservice.CookieService
}

func NewHandler(withdrawService withdrawservice.WithdrawService, cookieService cookieservice.CookieService) WithdrawHandler {
	return &withdrawHandlerImpl{withdrawService, cookieService}
}

func (w *withdrawHandlerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	userID := w.cookieService.AuthenticateUser(writer, request)

	var req entity.WithdrawRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	err = w.withdrawService.Create(request.Context(), req, userID)

	if err != nil {
		var le *myerrors.LimitExceededError
		var iof *myerrors.InvalidOrderNumberFormatError
		if errors.As(err, &le) {
			http.Error(writer, err.Error(), http.StatusPaymentRequired)
			return
		}

		if errors.As(err, &iof) {
			http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (w *withdrawHandlerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	//userID := o.cookieService.AuthenticateUser(writer, request)

	//ordersList, notFoundErr := o.orderService.GetAllByUserID(request.Context(), userID)
	//if notFoundErr != nil {
	//	http.Error(writer, notFoundErr.Error(), http.StatusNoContent)
	//	return
	//}
	//
	//writer.Header().Set("Content-Type", "application/json")
	//writer.WriteHeader(http.StatusOK)
	//writeErr := json.NewEncoder(writer).Encode(ordersList)
	//
	//if writeErr != nil {
	//	http.Error(writer, writeErr.Error(), http.StatusInternalServerError)
	//	return
	//}
}
