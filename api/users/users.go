package users

import (
	"github.com/Kaibling/psychic-octo-stock/middleware"
	"github.com/go-chi/chi"
)

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.Post("/", userPost)
	r.Get("/", usersGet)
	r.With(middleware.Authorization).Put("/{id}", userPut)
	r.With(middleware.Authorization).Delete("/{id}", userDelete)
	r.Get("/{id}", userGet)
	return r
}
