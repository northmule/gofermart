package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/northmule/gophermart/internal/app/api/handlers"
	"github.com/northmule/gophermart/internal/app/repository"
	"github.com/northmule/gophermart/internal/app/services/logger"
	"github.com/northmule/gophermart/internal/app/services/order"
	"net/http"
)

type AppRoutes struct {
	manager repository.Repository
}

func NewAppRoutes(repositoryManager repository.Repository) AppRoutes {
	instance := AppRoutes{
		manager: repositoryManager,
	}
	return instance
}

// DefiningAppRoutes маршруты приложения
func (ar *AppRoutes) DefiningAppRoutes(ctx context.Context) chi.Router {

	// Обработчики
	finalizeHandler := handlers.NewFinalizeHandler()
	registrationHandler := handlers.NewRegistrationHandler(ar.manager)
	checkAuthenticationHandler := handlers.NewCheckAuthenticationHandler(ar.manager)
	orderService := order.NewOrderService()
	orderHandler := handlers.NewOrderHandler(ar.manager, orderService)
	balanceHandler := handlers.NewBalanceHandler(ar.manager)
	withdrawHandler := handlers.NewWithdrawHandler(ar.manager, orderService)
	jobHandler := handlers.NewJobHandler(ar.manager)
	accrualHandler := handlers.NewAccrualHandler(ar.manager)

	r := chi.NewRouter()

	// Общие мидлвары
	r.Use(middleware.RequestLogger(logger.LogSugar))
	// r.Use(handlers.AddCommonContext(ctx))
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Route("/api/user", func(r chi.Router) {
		// регистрация пользователя
		r.With(
			handlers.AddCommonContext(ctx),
			registrationHandler.Registration,
			balanceHandler.CreateUserBalance,
			registrationHandler.Authentication,
		).Post("/register", finalizeHandler.FinalizeOk)

		// аутентификация пользователя
		r.With(
			handlers.AddCommonContext(ctx),
			registrationHandler.AuthenticationFromForm,
			registrationHandler.Authentication,
		).Post("/login", finalizeHandler.FinalizeOk)

		// загрузка пользователем номера заказа для расчёта
		r.With(
			handlers.AddCommonContext(ctx),
			checkAuthenticationHandler.Check,
			orderHandler.UploadingOrder,
			accrualHandler.CreateZeroAccrualForOrder,
			jobHandler.CreateJobToProcessNewOrder,
		).Post("/orders", func(res http.ResponseWriter, req *http.Request) {

		})

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
