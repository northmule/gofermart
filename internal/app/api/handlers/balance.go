package handlers

import (
	"github.com/northmule/gofermart/internal/app/api/rctx"
	"github.com/northmule/gofermart/internal/app/repository"
	"github.com/northmule/gofermart/internal/app/repository/models"
	"github.com/northmule/gofermart/internal/app/services/logger"
	"net/http"
)

type BalanceHandler struct {
	manager *repository.Manager
}

func NewBalanceHandler(manager *repository.Manager) *BalanceHandler {
	instance := &BalanceHandler{
		manager: manager,
	}
	return instance
}

type responseBalance struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

func (bh *BalanceHandler) Balance(res http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(rctx.UserCtxKey).(models.User)
	logger.LogSugar.Infof("Поступил запрос %s от пользователя %s", req.URL.Path, user.UUID)
}
