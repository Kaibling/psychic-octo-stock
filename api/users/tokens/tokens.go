package tokens

import (
	"github.com/Kaibling/psychic-octo-stock/apimiddleware"
	"github.com/go-chi/chi"
)

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.With(apimiddleware.Authorization).Get("/", tokensGet)
	r.With(apimiddleware.Authorization).Post("/", tokenPost)
	r.With(apimiddleware.Authorization).Put("/{tid}", tokenPut)
	r.With(apimiddleware.Authorization).Delete("/{tid}", tokenDelete)
	// r.With(apimiddleware.Authorization).Get("/{tid}", tokenGet)
	return r
}
