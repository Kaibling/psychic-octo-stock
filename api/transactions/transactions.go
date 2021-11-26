package transactions

import (
	"github.com/Kaibling/psychic-octo-stock/apimiddleware"
	"github.com/go-chi/chi"
)

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.With(apimiddleware.Authorization).Post("/", transactionPost)
	r.With(apimiddleware.Authorization).Get("/", transactionsGet)
	r.With(apimiddleware.Authorization).Delete("/{id}", transactionDelete)
	r.With(apimiddleware.Authorization).Get("/{id}", transactionGet)
	return r
}
