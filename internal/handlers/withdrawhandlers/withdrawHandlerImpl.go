package withdrawhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/security"
	"github.com/SemenRyzhkov/go-market-app/internal/service/withdrawservice"
)

type withdrawHandlerImpl struct {
	withdrawService withdrawservice.WithdrawService
	jwtHelper       *security.JwtHelper
}

func NewHandler(withdrawService withdrawservice.WithdrawService, jwtHelper *security.JwtHelper) WithdrawHandler {
	return &withdrawHandlerImpl{withdrawService, jwtHelper}
}

func (w *withdrawHandlerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	userID, err := w.jwtHelper.ExtractClaims(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
	var req entity.WithdrawRequest
	err = json.NewDecoder(request.Body).Decode(&req)
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

func (w *withdrawHandlerImpl) GetUserBalance(writer http.ResponseWriter, request *http.Request) {
	userID, err := w.jwtHelper.ExtractClaims(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	balanceRequest, err := w.withdrawService.GetUserBalance(request.Context(), userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	writeErr := json.NewEncoder(writer).Encode(balanceRequest)
	if writeErr != nil {
		http.Error(writer, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (w *withdrawHandlerImpl) GetAll(writer http.ResponseWriter, request *http.Request) {
	userID, err := w.jwtHelper.ExtractClaims(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	withdrawDTOList, notFoundErr := w.withdrawService.GetAllByUserID(request.Context(), userID)
	if notFoundErr != nil {
		http.Error(writer, notFoundErr.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	if len(withdrawDTOList) == 0 {
		writer.WriteHeader(http.StatusNoContent)
		return
	}
	writer.WriteHeader(http.StatusOK)

	writeErr := json.NewEncoder(writer).Encode(withdrawDTOList)
	if writeErr != nil {
		http.Error(writer, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}
