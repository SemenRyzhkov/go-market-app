package userhandlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
	"github.com/SemenRyzhkov/go-market-app/internal/entity/myerrors"
	"github.com/SemenRyzhkov/go-market-app/internal/security"
	"github.com/SemenRyzhkov/go-market-app/internal/service/userservice"
)

type userHandlerImpl struct {
	userService userservice.UserService
	jwtHelper   *security.JwtHelper
}

func NewHandler(userService userservice.UserService, jwtHelper *security.JwtHelper) UserHandler {
	return &userHandlerImpl{userService, jwtHelper}
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
	token, err := u.jwtHelper.GenerateJWT(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Authorization", token)
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

	token, err := u.jwtHelper.GenerateJWT(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	writer.Header().Set("Authorization", token)
}
