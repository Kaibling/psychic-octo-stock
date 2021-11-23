package authentication

import "github.com/go-chi/chi"

func AddRoute() chi.Router {
	r := chi.NewRouter()
	r.Post("/login", login)
	return r
}
