package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/northmule/gophermart/internal/app/api/handlers"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/services/accrual"
	"github.com/northmule/gophermart/internal/app/services/job"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/northmule/gophermart/internal/app/services/order"
	"github.com/northmule/gophermart/internal/app/storage"
)

type AppRoutes struct {
	manager repository.Repository
	storage storage.DBQuery
	session storage.SessionManager
}

func NewAppRoutes(repositoryManager repository.Repository, storage storage.DBQuery, session storage.SessionManager) AppRoutes {
	instance := AppRoutes{
		manager: repositoryManager,
		storage: storage,
		session: session,
	}
	return instance
}

// DefiningAppRoutes маршруты приложения
func (ar *AppRoutes) DefiningAppRoutes(ctx context.Context) chi.Router {

	//Сервисы
	orderService := order.NewOrderService()
	jobService := job.NewJobService(ar.manager)
	accrualService := accrual.NewAccrualService(ar.manager)

	// Обработчики
	registrationHandler := handlers.NewRegistrationHandler(ar.manager, ar.session)
	checkAuthenticationHandler := handlers.NewCheckAuthenticationHandler(ar.manager, ar.session)
	orderHandler := handlers.NewOrderHandler(ar.manager, orderService, jobService, accrualService)
	balanceHandler := handlers.NewBalanceHandler(ar.manager)
	withdrawHandler := handlers.NewWithdrawHandler(ar.manager, orderService)
	transactionHandler := handlers.NewTransactionHandler(ar.storage)

	r := chi.NewRouter()

	// Общие мидлвары
	r.Use(middleware.RequestLogger(logger.LogSugar))
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Route("/api/user", func(r chi.Router) {
		// регистрация пользователя
		r.With(
			handlers.AddCommonContext(ctx),
			transactionHandler.Transaction, // весь запрос в рамках транзации
		).Post("/register", registrationHandler.Registration)

		// аутентификация пользователя
		r.With(
			handlers.AddCommonContext(ctx),
		).Post("/login", registrationHandler.Authentication)

		// загрузка пользователем номера заказа для расчёта
		r.With(
			handlers.AddCommonContext(ctx),
			checkAuthenticationHandler.Check,
		).Post("/orders", orderHandler.UploadingOrder)

		// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
		r.With(
			handlers.AddCommonContext(ctx),
			checkAuthenticationHandler.Check,
		).Get("/orders", orderHandler.OrderList)

		// получение текущего баланса счёта баллов лояльности пользователя
		r.With(
			handlers.AddCommonContext(ctx),
			checkAuthenticationHandler.Check,
		).Get("/balance", balanceHandler.Balance)

		// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
		r.With(
			handlers.AddCommonContext(ctx),
			checkAuthenticationHandler.Check,
		).Post("/balance/withdraw", withdrawHandler.Withdraw)

		// получение информации о выводе средств с накопительного счёта пользователем
		r.With(
			handlers.AddCommonContext(ctx),
			checkAuthenticationHandler.Check,
		).Get("/withdrawals", withdrawHandler.WithdrawalsList)
	})

	return r
}
