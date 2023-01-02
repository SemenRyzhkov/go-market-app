package userhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/service/cookieservice"
	"github.com/SemenRyzhkov/go-market-app/internal/service/userservice"
)

type userHandlerImpl struct {
	userService   userservice.UserService
	cookieService cookieservice.CookieService
}

func NewHandler(userService userservice.UserService, cookieService cookieservice.CookieService) UserHandler {
	return &userHandlerImpl{userService, cookieService}
}

func (u *userHandlerImpl) Create(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var user entity.UserRequest
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}

	userID := uuid.New().String()
	err = u.userService.Create(request.Context(), user, userID)

	if err != nil {
		var ve *myerrors.UserViolationError
		if errors.As(err, &ve) {
			writer.WriteHeader(http.StatusConflict)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	cookieErr := u.cookieService.WriteSigned(writer, userID)
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (u *userHandlerImpl) Login(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	var user entity.UserRequest
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	userID, err := u.userService.Login(request.Context(), user)
	if err != nil {
		var ip *myerrors.InvalidPasswordError
		if errors.As(err, &ip) {
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	cookieErr := u.cookieService.WriteSigned(writer, userID)
	if cookieErr != nil {
		http.Error(writer, cookieErr.Error(), http.StatusInternalServerError)
		return
	}
}
