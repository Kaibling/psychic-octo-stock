package transactions

import (
	"github.com/Kaibling/psychic-octo-stock/middleware"
	"github.com/go-chi/chi"
)

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.With(middleware.Authorization).Post("/", transactionPost)
	r.With(middleware.Authorization).Get("/", transactionsGet)
	r.With(middleware.Authorization).Delete("/{id}", transactionDelete)
	r.With(middleware.Authorization).Get("/{id}", transactionGet)
	return r
}
