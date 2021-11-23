package stocks

import (
	"github.com/Kaibling/psychic-octo-stock/middleware"
	"github.com/go-chi/chi"
)

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.With(middleware.Authorization).Post("/users/{userid}", stockPost)
	r.With(middleware.Authorization).Get("/", stocksGet)
	r.With(middleware.Authorization).Put("/{id}", stockPut)
	r.With(middleware.Authorization).Delete("/{id}", stockDelete)
	r.With(middleware.Authorization).Get("/{id}", stockGet)
	return r
}
