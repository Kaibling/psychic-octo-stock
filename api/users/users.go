package users

import (
	"github.com/Kaibling/psychic-octo-stock/api/users/funds"
	"github.com/Kaibling/psychic-octo-stock/api/users/tokens"
	"github.com/Kaibling/psychic-octo-stock/apimiddleware"
	"github.com/go-chi/chi"
)

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.With(apimiddleware.Authorization).Post("/", userPost)
	r.With(apimiddleware.Authorization).Get("/", usersGet)
	r.With(apimiddleware.Authorization).Put("/{id}", userPut)
	r.With(apimiddleware.Authorization).Delete("/{id}", userDelete)
	r.With(apimiddleware.Authorization).Get("/{id}", userGet)
	r.Mount("/{id}/tokens", tokens.AddRoute())
	r.Mount("/{id}/funds", funds.AddRoute())
	return r
}
