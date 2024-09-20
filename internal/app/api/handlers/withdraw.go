package handlers

import (
	"encoding/json"
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/repository/models"
	"github.com/northmule/gophermart/internal/app/services/logger"
	orderService "github.com/northmule/gophermart/internal/app/services/order"
	"io"
	"net/http"
	"time"
)

type WithdrawHandler struct {
	manager      repository.Repository
	orderService *orderService.OrderService
}

func NewWithdrawHandler(manager repository.Repository, orderService *orderService.OrderService) *WithdrawHandler {
	instance := &WithdrawHandler{
		manager:      manager,
		orderService: orderService,
	}
	return instance
}

type requestWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type responseWithdrawals struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func (wh *WithdrawHandler) Withdraw(res http.ResponseWriter, req *http.Request) {
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer req.Body.Close()
	user := req.Context().Value(rctx.UserCtxKey).(models.User)
	logger.LogSugar.Infof("Данные запроса: %s", string(rawBody))

	var request requestWithdraw
	if err = json.Unmarshal(rawBody, &request); err != nil {
		logger.LogSugar.Infof("Пришли данные на списание %s. Запрос вызвал ошибку.", string(rawBody))
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if valid := wh.orderService.ValidateOrderNumber(request.Order); !valid {
		logger.LogSugar.Infof("Неверный формат номера заказа %s", request.Order)
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userBalance, err := wh.manager.Balance().FindOneByUserUUID(req.Context(), user.UUID)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	var balanceValue float64
	if userBalance != nil {
		balanceValue = userBalance.Value
	}

	if request.Sum > balanceValue {
		logger.LogSugar.Infof("Не достаточный баланс для списания по заказу %s. Текущий баланс:%f к списанию %f", request.Order, balanceValue, request.Sum)
		res.WriteHeader(http.StatusPaymentRequired)
		return
	}

	// order, err := wh.manager.Order().FindOneByNumber(req.Context(), request.Order)
	order, err := wh.manager.Order().FindByNumberOrCreate(req.Context(), request.Order, user.UUID)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if order == nil || order.ID == 0 {
		logger.LogSugar.Infof("Заказ с номером %s, по которому требуется списать %f ещё не создан.", request.Order, request.Sum)
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	existWithdraw, err := wh.manager.Withdrawn().FindOneByOrderID(req.Context(), order.ID)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if existWithdraw != nil && existWithdraw.ID > 0 {
		logger.LogSugar.Infof("Списание по заказу уже было выполнено для заказа с номером %s", request.Order)
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	_, err = wh.manager.Withdrawn().Withdraw(req.Context(), user.UUID, request.Sum, order.ID)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.LogSugar.Infof("Списание на сумуу %f удачно выполнено для заказа %s", request.Sum, order.Number)
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func (wh *WithdrawHandler) WithdrawalsList(res http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(rctx.UserCtxKey).(models.User)

	withdraws, err := wh.manager.Withdrawn().FindWithdrawsByUserUUID(req.Context(), user.UUID)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdraws) == 0 {
		logger.LogSugar.Infof("Нет данных по списаниям для пользователя с uuuid %s", user.UUID)
		res.WriteHeader(http.StatusNoContent)
		return
	}

	var responseList []responseWithdrawals

	for _, withdrawn := range withdraws {
		response := responseWithdrawals{
			Order:       withdrawn.Order.Number,
			Sum:         withdrawn.Value,
			ProcessedAt: withdrawn.CreatedAt.Format(time.RFC3339),
		}
		responseList = append(responseList, response)
	}

	responseListValue, err := json.Marshal(responseList)
	if err != nil {
		logger.LogSugar.Errorf(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseListValue)
	if err != nil {
		http.Error(res, "Ответ не передан", http.StatusInternalServerError)
		return
	}
}
