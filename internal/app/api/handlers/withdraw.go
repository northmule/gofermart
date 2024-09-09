package handlers

import (
	"encoding/json"
	"github.com/northmule/gofermart/internal/app/api/rctx"
	"github.com/northmule/gofermart/internal/app/repository"
	"github.com/northmule/gofermart/internal/app/repository/models"
	"github.com/northmule/gofermart/internal/app/services/logger"
	orderService "github.com/northmule/gofermart/internal/app/services/order"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

type WithdrawHandler struct {
	manager          *repository.Manager
	orderService     *orderService.OrderService
	regexOrderNumber *regexp.Regexp
}

func NewWithdrawHandler(manager *repository.Manager, orderService *orderService.OrderService) *WithdrawHandler {
	instance := &WithdrawHandler{
		manager:          manager,
		orderService:     orderService,
		regexOrderNumber: regexp.MustCompile(`\d+`),
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
	logger.LogSugar.Infof("Поступил запрос %s от пользователя %s. Данные запроса: %s", req.URL.Path, user.UUID, string(rawBody))

	var request requestWithdraw
	if err = json.Unmarshal(rawBody, &request); err != nil {
		logger.LogSugar.Infof("Пришли данные на списание %s. Запрос вызвал ошибку.", string(rawBody))
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if !wh.regexOrderNumber.MatchString(request.Order) || !wh.validateOrderNumber(request.Order) {
		logger.LogSugar.Infof("Неверный формат номера заказа %s", request.Order)
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	userBalance, err := wh.manager.Balance.FindOneByUserUUID(user.UUID)
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

	order, err := wh.manager.Order.FindOneByNumber(request.Order)
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
	existWithdraw, err := wh.manager.Withdrawn.FindOneByOrderID(order.ID)
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
	_, err = wh.manager.Withdrawn.Withdraw(user.ID, request.Sum, order.ID)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	logger.LogSugar.Infof("Списание на сумуу %f удачно выполнено для заказа %s", request.Sum, order.Number)
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
}

func (wh *WithdrawHandler) validateOrderNumber(orderNumber string) bool {
	orderInt, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		return false
	}
	return wh.orderService.ValidateOrderNumber(int(orderInt))
}

func (wh *WithdrawHandler) WithdrawalsList(res http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(rctx.UserCtxKey).(models.User)
	logger.LogSugar.Infof("Поступил запрос %s от пользователя %s", req.URL.Path, user.UUID)

	withdraws, err := wh.manager.Withdrawn.FindWithdrawsByUserUUID(user.UUID)
	if err != nil {
		logger.LogSugar.Error(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(*withdraws) == 0 {
		logger.LogSugar.Infof("Нет данных по списаниям для пользователя с uuuid %s", user.UUID)
		res.WriteHeader(http.StatusNoContent)
		return
	}

	var responseList []responseWithdrawals

	for _, withdrawn := range *withdraws {
		response := responseWithdrawals{
			Order:       withdrawn.Order.Number,
			Sum:         withdrawn.Value,
			ProcessedAt: time.Unix(withdrawn.CreatedAt.Time.Unix(), 0).Format(time.RFC3339),
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
