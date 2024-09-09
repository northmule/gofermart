package handlers

import (
	"github.com/northmule/gofermart/internal/app/api/rctx"
	"github.com/northmule/gofermart/internal/app/repository"
	"github.com/northmule/gofermart/internal/app/repository/models"
	"github.com/northmule/gofermart/internal/app/services/logger"
	"net/http"
)

type JobHandler struct {
	manager *repository.Manager
}

func NewJobHandler(manager *repository.Manager) *JobHandler {
	instance := &JobHandler{
		manager: manager,
	}
	return instance
}

func (jh *JobHandler) CreateTaskToProcessNewOrder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		newOrder := req.Context().Value(rctx.OrderUpload).(models.Order)
		_, err := jh.manager.Job.CreateJobByOrderNumber(newOrder.Number)
		if err != nil {
			logger.LogSugar.Errorf("Ошибка создания задания на обработку заказа с номером %s", newOrder.Number)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(res, req)
	})
}
