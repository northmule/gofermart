package handlers

import (
	"context"
	"encoding/json"
	"github.com/northmule/gofermart/internal/app/api/rctx"
	"github.com/northmule/gofermart/internal/app/constants"
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

type OrderHandler struct {
	manager          *repository.Manager
	orderService     *orderService.OrderService
	regexOrderNumber *regexp.Regexp
}

type orderResponse struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

func NewOrderHandler(manager *repository.Manager, orderService *orderService.OrderService) *OrderHandler {
	instance := &OrderHandler{
		manager:          manager,
		orderService:     orderService,
		regexOrderNumber: regexp.MustCompile(`\d+`),
	}
	return instance
}

func (o *OrderHandler) UploadingOrder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		rawBody, err := io.ReadAll(req.Body)
		if err != nil {
			logger.LogSugar.Errorf(err.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer req.Body.Close()
		if len(rawBody) == 0 {
			logger.LogSugar.Infof("Пустое тело запроса. Тело запроса должно содержать номер заказа.")
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		orderNumber := string(rawBody)
		if !o.regexOrderNumber.MatchString(orderNumber) || !o.validateOrderNumber(orderNumber) {
			logger.LogSugar.Infof("Неверный формат номера заказа %s", orderNumber)
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		user := req.Context().Value(rctx.UserCtxKey).(models.User)
		logger.LogSugar.Infof("Поступил запрос %s от пользователя %s", req.URL.Path, user.UUID)
		logger.LogSugar.Infof("Получин номер заказа %s, от пользователя %s", orderNumber, user.Login)

		order, err := o.manager.Order.FindOneByNumber(orderNumber)
		if err != nil {
			logger.LogSugar.Errorf(err.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Заказ найден по номеру
		if order != nil && order.ID > 0 && order.User.ID > 0 {
			if order.User.UUID == user.UUID {
				logger.LogSugar.Infof("Заказ %s уже был загружен текущим пользователем", orderNumber)
				res.WriteHeader(http.StatusOK)
				return
			}
			if order.User.UUID != user.UUID {
				logger.LogSugar.Infof("Заказ %s уже был загружен другим пользователем: %s", orderNumber, order.User.Login)
				res.WriteHeader(http.StatusConflict)
				return
			}
			logger.LogSugar.Errorf("Не ожиданное поведение: %v -  %v", order, user)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		newOrder := models.Order{
			Number: orderNumber,
			Status: constants.OrderStatusNew,
			User:   user,
		}
		orderID, err := o.manager.Order.Save(newOrder, newOrder.User.ID)
		if err != nil {
			logger.LogSugar.Errorf(err.Error())
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if orderID == 0 {
			logger.LogSugar.Error("Не присвоен ID заказа после сохранения")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		newOrder.ID = int(orderID)
		res.WriteHeader(http.StatusAccepted)
		ctx := context.WithValue(req.Context(), rctx.OrderUpload, newOrder)
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)
	})
}

func (o *OrderHandler) OrderList(res http.ResponseWriter, req *http.Request) {
	user := req.Context().Value(rctx.UserCtxKey).(models.User)
	logger.LogSugar.Infof("Поступил запрос %s от пользователя %s", req.URL.Path, user.UUID)

	orders, err := o.manager.Order.FindOrdersByUserUUID(user.UUID)
	if err != nil {
		logger.LogSugar.Errorf(err.Error())
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	if orders == nil || len(*orders) == 0 {
		logger.LogSugar.Infof("Нет заказов для отображения для пользователя %s", user.UUID)
		res.WriteHeader(http.StatusNoContent)
		return
	}
	var orderListResponse []orderResponse
	for _, order := range *orders {
		orderResponse := orderResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual.Float64,
			UploadedAt: time.Unix(order.CreatedAt.Time.Unix(), 0).Format(time.RFC3339),
		}
		orderListResponse = append(orderListResponse, orderResponse)
	}

	orderListResponseValue, err := json.Marshal(orderListResponse)
	if err != nil {
		http.Error(res, "Ошибка подготовки ответа", http.StatusInternalServerError)
		logger.LogSugar.Errorf(err.Error())
		return
	}
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(orderListResponseValue)
	if err != nil {
		http.Error(res, "Ответ не передан", http.StatusInternalServerError)
		return
	}
}

func (o *OrderHandler) validateOrderNumber(orderNumber string) bool {
	orderInt, err := strconv.ParseInt(orderNumber, 10, 64)
	if err != nil {
		logger.LogSugar.Errorf(err.Error())
		return false
	}
	return o.orderService.ValidateOrderNumber(int(orderInt))
}
