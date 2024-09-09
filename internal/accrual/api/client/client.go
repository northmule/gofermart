package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/northmule/gofermart/config"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

const ServiceAccrualURI = "/api/orders/{number}"
const ServiceName = "Accrual"

type ErrorNoContent string
type ErrorTooManyRequests string
type ErrorInternalServerError string
type ErrorUndefined struct {
	code int
}

func (err ErrorNoContent) Error() string {
	return fmt.Sprintf("Сервис %s вернул пустой ответ", ServiceName)
}

func (err ErrorTooManyRequests) Error() string {
	return fmt.Sprintf("Сервис %s перегружен запросами", ServiceName)
}

func (err ErrorInternalServerError) Error() string {
	return fmt.Sprintf("Сервис %s сломался", ServiceName)
}

func (err ErrorUndefined) Error() string {
	return fmt.Sprintf("Сервис %s вернул не обработаный код ошибки %d", ServiceName, err.code)
}

type AccrualClient struct {
	serviceURL string
	logger     *zap.SugaredLogger
}

func NewAccrualClient(serviceURL string, logger *zap.SugaredLogger) *AccrualClient {
	instance := &AccrualClient{
		serviceURL: serviceURL,
		logger:     logger,
	}
	return instance
}

type ResponseAccrual struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (ac *AccrualClient) SendOrderNumber(orderNumber string) (*ResponseAccrual, error) {

	ctx, cancel := context.WithTimeout(context.Background(), config.DataBaseConnectionTimeOut*time.Second)
	defer cancel()
	requestURL := fmt.Sprintf(
		"%s%s", strings.TrimRight(ac.serviceURL, "/"),
		strings.Replace(ServiceAccrualURI, "{number}", orderNumber, 1),
	)
	ac.logger.Infof("Поступил запрос: %s на полученине информации от сервиса %s", requestURL, ServiceName)
	requestPrepare, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	accrualResponse, err := client.Do(requestPrepare)
	if err != nil {
		return nil, err
	}
	defer accrualResponse.Body.Close()
	if ok, err := ac.isStatusOk(accrualResponse); !ok {
		return nil, err
	}

	responseBodyRaw, err := io.ReadAll(accrualResponse.Body)
	if err != nil {
		return nil, err
	}
	ac.logger.Infof("Получен ответ от сервиса %s: %s", ServiceName, string(responseBodyRaw))

	response := &ResponseAccrual{}

	err = json.Unmarshal(responseBodyRaw, response)
	if err != nil {
		ac.logger.Infof("Не удалось разобрать ответ от севриса %s", ServiceName)
		return nil, err
	}

	return response, nil
}

func (ac *AccrualClient) isStatusOk(response *http.Response) (bool, error) {
	if response.StatusCode == http.StatusOK {
		return true, nil
	}
	if response.StatusCode == http.StatusNoContent {
		return false, new(ErrorNoContent)
	}

	if response.StatusCode == http.StatusTooManyRequests {
		return false, new(ErrorTooManyRequests)
	}

	if response.StatusCode == http.StatusInternalServerError {
		return false, new(ErrorInternalServerError)
	}

	return false, new(ErrorUndefined)
}
