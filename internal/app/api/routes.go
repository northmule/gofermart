package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/northmule/gofermart/internal/app/api/handlers"
	"github.com/northmule/gofermart/internal/app/repository"
	"github.com/northmule/gofermart/internal/app/services/order"
	"net/http"
)

type AppRoutes struct {
	manager *repository.Manager
}

func NewAppRoutes(repositoryManager *repository.Manager) chi.Router {
	instance := AppRoutes{
		manager: repositoryManager,
	}
	return instance.DefiningAppRoutes()
}

// DefiningAppRoutes маршруты приложения
func (ar *AppRoutes) DefiningAppRoutes() chi.Router {
	r := chi.NewRouter()

	finalizeHandler := handlers.NewFinalizeHandler()
	registrationHandler := handlers.NewRegistrationHandler(ar.manager)
	checkAuthenticationHandler := handlers.NewCheckAuthenticationHandler(ar.manager)

	orderService := order.NewOrderService()
	orderHandler := handlers.NewOrderHandler(ar.manager, orderService)

	balanceHandler := handlers.NewBalanceHandler(ar.manager)

	withdrawHandler := handlers.NewWithdrawHandler(ar.manager, orderService)

	jobHandler := handlers.NewJobHandler(ar.manager)

	accrualHandler := handlers.NewAccrualHandler(ar.manager)

	stub := func(res http.ResponseWriter, req *http.Request) {
		return
	}

	r.Route("/api/user", func(r chi.Router) {
		// регистрация пользователя
		r.With(
			registrationHandler.Registration,
			balanceHandler.CreateUserBalance,
			registrationHandler.Authentication,
		).Post("/register", finalizeHandler.FinalizeOk)

		// аутентификация пользователя
		r.With(
			registrationHandler.AuthenticationFromForm,
			registrationHandler.Authentication,
		).Post("/login", finalizeHandler.FinalizeOk)

		// загрузка пользователем номера заказа для расчёта
		r.With(
			checkAuthenticationHandler.Check,
			orderHandler.UploadingOrder,
			accrualHandler.CreateZeroAccrualForOrder,
			jobHandler.CreateTaskToProcessNewOrder,
		).Post("/orders", stub)

		// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
		r.With(
			checkAuthenticationHandler.Check,
		).Get("/orders", orderHandler.OrderList)

		// получение текущего баланса счёта баллов лояльности пользователя
		r.With(
			checkAuthenticationHandler.Check,
		).Get("/balance", balanceHandler.Balance)

		// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
		r.With(
			checkAuthenticationHandler.Check,
		).Post("/balance/withdraw", withdrawHandler.Withdraw)

		// получение информации о выводе средств с накопительного счёта пользователем
		r.With(
			checkAuthenticationHandler.Check,
		).Get("/withdrawals", withdrawHandler.WithdrawalsList)
	})

	return r
}
