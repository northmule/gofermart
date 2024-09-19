package handlers

import (
	"encoding/json"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"net/http"
)

type FinalizeHandler struct {
}

func NewFinalizeHandler() *FinalizeHandler {
	instance := &FinalizeHandler{}
	return instance
}

type Response struct {
	Ok bool `json:"ok"`
}

func (handler *FinalizeHandler) FinalizeOk(res http.ResponseWriter, req *http.Request) {
	logger.LogSugar.Info("Обработка запроса завершена")
	var response Response
	response.Ok = true
	responseValue, err := json.Marshal(response)
	if err != nil {
		http.Error(res, "Ошибка подготовки ответа", http.StatusInternalServerError)
		logger.LogSugar.Error(err)
		return
	}
	res.Header().Set("content-type", "application/json")
	res.WriteHeader(http.StatusOK)
	_, err = res.Write(responseValue)
	if err != nil {
		http.Error(res, "Ответ не передан", http.StatusInternalServerError)
		return
	}
}
