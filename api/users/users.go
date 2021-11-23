package users

import (
	"github.com/Kaibling/psychic-octo-stock/middleware"
	"github.com/go-chi/chi"
)

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.With(middleware.Authorization).Post("/", userPost)
	r.With(middleware.Authorization).Get("/", usersGet)
	r.With(middleware.Authorization).Put("/{id}", userPut)
	r.With(middleware.Authorization).Delete("/{id}", userDelete)
	r.With(middleware.Authorization).Get("/{id}", userGet)
	return r
}
