package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/northmule/gofermart/internal/app/repository"
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
	})

	return r
}
