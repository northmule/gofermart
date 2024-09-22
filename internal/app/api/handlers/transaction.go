package handlers

import (
	"context"
	"github.com/northmule/gophermart/internal/app/api/rctx"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/northmule/gophermart/internal/app/storage"
	"net/http"
)

type TransactionHandler struct {
	db storage.DBQuery
}

func NewTransactionHandler(db storage.DBQuery) *TransactionHandler {
	instance := &TransactionHandler{
		db: db,
	}

	return instance
}

func (th *TransactionHandler) Transaction(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var err error

		transaction, err := storage.NewTransaction(th.db)
		logger.LogSugar.Info("Открыл транзакцию")
		if err != nil {
			logger.LogSugar.Errorf("Транзакция не открыта: %s", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(req.Context(), rctx.TransactionCtxKey, transaction)
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)

		if len(transaction.Error()) > 0 {
			logger.LogSugar.Info("Во время выполнения транзакции, возникли ошибки", transaction.Error())
			if err = transaction.Rollback(); err != nil {
				logger.LogSugar.Errorf("Ошибка commit запроса: %s", err)

			}
			logger.LogSugar.Info("Произвёл откат транзакции, данные не изменены")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = transaction.Commit()
		if err != nil {
			logger.LogSugar.Errorf("Ошибка commit запроса: %s", err)
			res.WriteHeader(http.StatusInternalServerError)
		}
		logger.LogSugar.Info("Транзакция завершена")
	})
}
