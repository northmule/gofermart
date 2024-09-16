package handlers

import (
	"context"
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/services/authentication"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"net/http"
	"strings"
)

type CheckAuthenticationHandler struct {
	manager repository.Repository
}

func NewCheckAuthenticationHandler(manager repository.Repository) *CheckAuthenticationHandler {
	instance := &CheckAuthenticationHandler{
		manager: manager,
	}
	return instance
}

func (c *CheckAuthenticationHandler) Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		tokenValue := authentication.GetUserToken(req)
		if tokenValue == "" {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		cookieValues := strings.Split(tokenValue, ":")
		if len(cookieValues) < 2 {
			logger.LogSugar.Infof("UUID пользователя в куке не найден %s", tokenValue)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := cookieValues[0]
		userUUID := cookieValues[1]
		logger.LogSugar.Infof("Проверка UUID %s из значений cookie", userUUID)
		if !authentication.ValidateToken(userUUID, token, authentication.HMACSecretKey) {
			logger.LogSugar.Infof("Значение UUID %s из cookie не прошло проверку подписи", userUUID)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := c.manager.User().FindOneByUUID(userUUID)
		if err != nil {
			logger.LogSugar.Errorf("Ошибка при поиске пользователя по UUID %s, %s", userUUID, err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if user == nil || user.ID == 0 {
			logger.LogSugar.Infof("Пользователь с uuid %s не найден", userUUID)
			res.WriteHeader(http.StatusUnauthorized)
			return
		}
		logger.LogSugar.Infof("Поступил запрос %s от пользователя %s", req.URL.Path, user.UUID)

		ctx := context.WithValue(req.Context(), rctx.UserCtxKey, *user)
		req = req.WithContext(ctx)
		next.ServeHTTP(res, req)
	})
}
