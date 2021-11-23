package api

import (
	"github.com/Kaibling/psychic-octo-stock/api/authentication"
	"github.com/Kaibling/psychic-octo-stock/api/stocks"
	"github.com/Kaibling/psychic-octo-stock/api/transactions"
	"github.com/Kaibling/psychic-octo-stock/api/users"
	"github.com/go-chi/chi"
)

func BuildRouter(r *chi.Mux) *chi.Mux {

	r.Route("/api/v1", func(r chi.Router) {
		r.Mount("/users", users.AddRoute())
		r.Mount("/stocks", stocks.AddRoute())
		r.Mount("/transactions", transactions.AddRoute())
		r.Mount("/", authentication.AddRoute())
	})

	return r
}
