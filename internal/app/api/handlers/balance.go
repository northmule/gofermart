package handlers

import (
	"encoding/json"
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"net/http"
)

type BalanceHandler struct {
	manager repository.Repository
}

func NewBalanceHandler(manager repository.Repository) *BalanceHandler {
	instance := &BalanceHandler{
		manager: manager,
	}
	return instance
}

type responseBalance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

func (bh *BalanceHandler) Balance(res http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(rctx.UserCtxKey).(models.User)

	balance, err := bh.manager.Balance().FindOneByUserUUID(user.UUID)
	if err != nil {
		logger.LogSugar.Errorf(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	withdrawn, err := bh.manager.Withdrawn().FindSumWithdrawnByUserUUID(user.UUID)
	if err != nil {
		logger.LogSugar.Errorf(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := responseBalance{
		Current:   0,
		Withdrawn: withdrawn,
	}
	if balance != nil {
		response.Current = balance.Value
	}

	responseBalanceValue, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Ошибка подготовки ответа", http.StatusInternalServerError)
		logger.LogSugar.Error(err.Error())
		return
	}
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseBalanceValue)
	if err != nil {
		http.Error(res, "Ответ не передан", http.StatusInternalServerError)
		return
	}
}

func (bh *BalanceHandler) CreateUserBalance(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		newUser := req.Context().Value(rctx.UserCtxKey).(models.User)
		_, err := bh.manager.Balance().CreateBalanceByUserUUID(newUser.UUID)
		if err != nil {
			logger.LogSugar.Errorf(err.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(res, req)
	})
}
