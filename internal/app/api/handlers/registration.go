package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/northmule/gophermart/internal/app/services/authentication"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/northmule/gophermart/internal/app/storage"
	"github.com/northmule/gophermart/internal/app/util"
	"io"
	"net/http"
)

type RegistrationHandler struct {
	manager repository.Repository
	session storage.SessionManager
}

type registrationRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authenticationRequestBody struct {
	registrationRequestBody
}

func NewRegistrationHandler(manager repository.Repository, session storage.SessionManager) *RegistrationHandler {
	instance := &RegistrationHandler{
		manager: manager,
		session: session,
	}
	return instance
}

func (r *RegistrationHandler) Registration(res http.ResponseWriter, req *http.Request) {

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.LogSugar.Error(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	var request registrationRequestBody
	if err = json.Unmarshal(rawBody, &request); err != nil {
		logger.LogSugar.Infof("Пришли данные регистрации: %s. Запрос вызвал ошибку.", string(rawBody))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := r.manager.User().FindOneByLogin(req.Context(), request.Login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user != nil && user.Login == request.Login {
		logger.LogSugar.Infof("Логин '%s' уже занят  пользователем с uuid '%s'", request.Login, user.UUID)
		res.WriteHeader(http.StatusConflict)
		return
	}

	newUser := models.User{
		Login:    request.Login,
		Password: util.PasswordHash(request.Password),
		UUID:     uuid.NewString(),
	}
	tx := req.Context().Value(rctx.TransactionCtxKey).(*storage.Transaction)
	userID, err := r.manager.User().TxCreateNewUser(req.Context(), tx, newUser)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userID == 0 {
		logger.LogSugar.Error("Пустое значение ID при регистрации пользователя")
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser.ID = int(userID)
	logger.LogSugar.Infof("Зарегистрирован новый пользователь '%s' с uuid '%s'", newUser.Login, newUser.UUID)

	ctx := context.WithValue(req.Context(), rctx.UserCtxKey, newUser)
	req = req.WithContext(ctx)

	_, err = r.manager.Balance().TxCreateBalanceByUserUUID(req.Context(), tx, newUser.UUID)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, tokenValue, cookie, tokenExp := authentication.Authentication(newUser.UUID)

	http.SetCookie(res, cookie)
	r.session.Add(token, *tokenExp)

	res.Header().Set("Authorization", tokenValue)
	logger.LogSugar.Infof("Пользователь прошёл аунтификацию, выдан токен с uuid %s", tokenValue)
	res.WriteHeader(http.StatusOK)
}

func (r *RegistrationHandler) Authentication(res http.ResponseWriter, req *http.Request) {

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.LogSugar.Errorf(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()

	var request authenticationRequestBody
	if err = json.Unmarshal(rawBody, &request); err != nil {
		logger.LogSugar.Infof("Пришли данные аунтификации: %s. Запрос вызвал ошибку.", string(rawBody))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := r.manager.User().FindOneByLogin(req.Context(), request.Login)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if user == nil || len(request.Password) == 0 || user.Password != util.PasswordHash(request.Password) {
		logger.LogSugar.Infof("Неверная пара логин/пароль %s/****", request.Login)
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	logger.LogSugar.Infof("Аунтификация пользователя: %s", user.UUID)
	if user.UUID == "" {
		logger.LogSugar.Info("Пользователь не распознан для аунтификации")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, tokenValue, cookie, tokenExp := authentication.Authentication(user.UUID)

	http.SetCookie(res, cookie)
	r.session.Add(token, *tokenExp)

	res.Header().Set("Authorization", tokenValue)
	logger.LogSugar.Infof("Пользователь прошёл аунтификацию, выдан токен с uuid %s", tokenValue)
	res.WriteHeader(http.StatusOK)

}
