package orderhandlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/security"
	"github.com/SemenRyzhkov/go-market-app/internal/service/orderservice"
)

type orderHandlerImpl struct {
	orderService orderservice.OrderService
	jwtHelper    *security.JwtHelper
}

func NewHandler(orderService orderservice.OrderService, jwtHelper *security.JwtHelper) OrderHandler {
	return &orderHandlerImpl{orderService, jwtHelper}
}

func (o *orderHandlerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := o.jwtHelper.ExtractClaims(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	orderNumber, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	err = o.orderService.Create(request.Context(), string(orderNumber), userID)

	if err != nil {
		var ov *myerrors.OrderViolationError
		var eo *myerrors.ExistedOrderError
		var iof *myerrors.InvalidOrderNumberFormatError
		if errors.As(err, &ov) {
			http.Error(writer, err.Error(), http.StatusOK)
			return
		}
		if errors.As(err, &eo) {
			http.Error(writer, err.Error(), http.StatusConflict)
			return
		}
		if errors.As(err, &iof) {
			http.Error(writer, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusAccepted)
}

func (o *orderHandlerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := o.jwtHelper.ExtractClaims(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	ordersList, notFoundErr := o.orderService.GetAllByUserID(request.Context(), userID)
	if notFoundErr != nil {
		http.Error(writer, notFoundErr.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writeErr := json.NewEncoder(writer).Encode(ordersList)

	if writeErr != nil {
		http.Error(writer, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}
