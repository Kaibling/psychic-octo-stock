package stocks

import (
	"github.com/Kaibling/psychic-octo-stock/apimiddleware"
	"github.com/go-chi/chi"
)

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.With(apimiddleware.Authorization).Post("/users/{userid}", stockPost)
	r.With(apimiddleware.Authorization).Get("/", stocksGet)
	r.With(apimiddleware.Authorization).Put("/{id}", stockPut)
	r.With(apimiddleware.Authorization).Delete("/{id}", stockDelete)
	r.With(apimiddleware.Authorization).Get("/{id}", stockGet)
	return r
}
