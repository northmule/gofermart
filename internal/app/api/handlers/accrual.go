package handlers

import (
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"go.uber.org/zap"
	"net/http"
)

type AccrualHandler struct {
	manager repository.Repository
}

func NewAccrualHandler(manager repository.Repository) *AccrualHandler {
	instance := &AccrualHandler{
		manager: manager,
	}
	return instance
}

func (bh *AccrualHandler) CreateZeroAccrualForOrder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		newUser := req.Context().Value(rctx.UserCtxKey).(models.User)
		newOrder := req.Context().Value(rctx.OrderUpload).(models.Order)
		_, err := bh.manager.Accrual().CreateAccrualByOrderNumberAndUserUUID(req.Context(), newOrder.Number, newUser.UUID)
		logger.LogSugar.Info("Создаю информацию о нулевом списании по заказу", zap.String("number", newOrder.Number))
		if err != nil {
			logger.LogSugar.Error(err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(res, req)
	})
}
