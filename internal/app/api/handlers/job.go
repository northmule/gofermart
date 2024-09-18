package handlers

import (
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"net/http"
)

type JobHandler struct {
	manager repository.Repository
}

func NewJobHandler(manager repository.Repository) *JobHandler {
	instance := &JobHandler{
		manager: manager,
	}
	return instance
}

func (jh *JobHandler) CreateTaskToProcessNewOrder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		newOrder := req.Context().Value(rctx.OrderUpload).(models.Order)
		_, err := jh.manager.Job().CreateJobByOrderNumber(req.Context(), newOrder.Number)

		if err != nil {
			logger.LogSugar.Errorf("Ошибка создания задания на обработку заказа с номером %s", newOrder.Number)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		logger.LogSugar.Infof("Создал задачу для запроса начисленных балов для заказа %s", newOrder.Number)
		next.ServeHTTP(res, req)
	})
}
