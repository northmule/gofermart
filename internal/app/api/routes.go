package api

import (
	"github.com/go-chi/chi/v5"
<<<<<<< HEAD
	"github.com/northmule/gofermart/internal/app/api/handlers"
	"github.com/northmule/gofermart/internal/app/repository"
	"github.com/northmule/gofermart/internal/app/services/order"
=======
	"github.com/northmule/gofermart/internal/app/repository"
	"net/http"
>>>>>>> 94746e2 (базовая структура)
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

<<<<<<< HEAD
	finalizeHandler := handlers.NewFinalizeHandler()
	registrationHandler := handlers.NewRegistrationHandler(ar.manager)
	checkAuthenticationHandler := handlers.NewCheckAuthenticationHandler(ar.manager)

	orderService := order.NewOrderService()
	orderHandler := handlers.NewOrderHandler(ar.manager, orderService)

	balanceHandler := handlers.NewBalanceHandler(ar.manager)

	withdrawHandler := handlers.NewWithdrawHandler(ar.manager, orderService)

	r.Route("/api/user", func(r chi.Router) {
		// регистрация пользователя
		r.With(
			registrationHandler.Registration,
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
		).Post("/orders", orderHandler.UploadingOrder)

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
=======
	mstub := func(next http.Handler) http.Handler {
		return next
	}

	r.Use(mstub)

	rstub := func(res http.ResponseWriter, req *http.Request) {

	}
	r.Route("/api/user", func(r chi.Router) {
		// регистрация пользователя
		r.Post("/register", rstub)
		// аутентификация пользователя
		r.Post("/login", rstub)
		// загрузка пользователем номера заказа для расчёта
		r.Post("/orders", rstub)
		// получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
		r.Get("/orders", rstub)
		// получение текущего баланса счёта баллов лояльности пользователя
		r.Get("/balance", rstub)
		// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
		r.Post("/balance/withdraw", rstub)
		// получение информации о выводе средств с накопительного счёта пользователем
		r.Get("/withdrawals", rstub)
>>>>>>> 94746e2 (базовая структура)
	})

	return r
}
