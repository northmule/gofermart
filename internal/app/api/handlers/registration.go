package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/northmule/gofermart/internal/app/api/rctx"
	"github.com/northmule/gofermart/internal/app/repository"
	"github.com/northmule/gofermart/internal/app/repository/models"
	"github.com/northmule/gofermart/internal/app/services/authentication"
	"github.com/northmule/gofermart/internal/app/services/logger"
	"github.com/northmule/gofermart/internal/app/util"
	"io"
	"net/http"
)

type RegistrationHandler struct {
	manager *repository.Manager
}

type registrationRequestBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type authenticationRequestBody struct {
	registrationRequestBody
}

func NewRegistrationHandler(manager *repository.Manager) *RegistrationHandler {
	instance := &RegistrationHandler{
		manager: manager,
	}
	return instance
}

func (r *RegistrationHandler) Registration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		rawBody, err := io.ReadAll(req.Body)
		if err != nil {
			logger.LogSugar.Errorf(err.Error())
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

		user, err := r.manager.User.FindOneByLogin(request.Login)
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
			Name:     "Имя",
			Login:    request.Login,
			Password: util.PasswordHash(request.Password),
			UUID:     uuid.NewString(),
		}

		userID, err := r.manager.User.CreateNewUser(newUser)
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

		next.ServeHTTP(res, req)
	})
}

func (r *RegistrationHandler) AuthenticationFromForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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

		user, err := r.manager.User.FindOneByLogin(request.Login)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if user == nil || len(request.Password) == 0 || user.Password != util.PasswordHash(request.Password) {
			logger.LogSugar.Infof("Неверная пара логин/пароль %s/****", request.Login)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(req.Context(), rctx.UserCtxKey, *user)
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)
	})
}

func (r *RegistrationHandler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		user := req.Context().Value(rctx.UserCtxKey).(models.User)
		logger.LogSugar.Infof("Аунтификация пользователя: %s", user.UUID)
		if user.UUID == "" {
			logger.LogSugar.Info("Пользователь не распознан для аунтификации")
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, tokenExp := authentication.GenerateToken(user.UUID, authentication.HMACTokenExp, authentication.HMACSecretKey)
		tokenValue := fmt.Sprintf("%s:%s", token, user.UUID)

		http.SetCookie(res, &http.Cookie{
			Name:    authentication.CookieAuthName,
			Value:   tokenValue,
			Expires: tokenExp,
			Secure:  false,
			Path:    "/",
		})

		res.Header().Set("Authorization", tokenValue)
		logger.LogSugar.Infof("Пользователь прошёл аунтификацию, выдан токен с uuid %s", tokenValue)

		next.ServeHTTP(res, req)
	})
}
