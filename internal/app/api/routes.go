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

	rstub := func(res http.ResponseWriter, req *http.Request) {
		return
	}
	finalizeHandler := handlers.NewFinalizeHandler()
	registrationHandler := handlers.NewRegistrationHandler(ar.manager)
	checkAuthenticationHandler := handlers.NewCheckAuthenticationHandler(ar.manager)

	orderService := order.NewOrderService()
	orderHandler := handlers.NewOrderHandler(ar.manager, orderService)

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
		r.Get("/balance", rstub)
		// запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа
		r.Post("/balance/withdraw", rstub)
		// получение информации о выводе средств с накопительного счёта пользователем
		r.Get("/withdrawals", rstub)
	})

	return r
}
