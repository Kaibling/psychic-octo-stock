package apimiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/repositories"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := transmission.GetResponse(r)

		auth := r.Header.Get("Authorization")
		if auth == "" {
			response.Send("", "Could not find Authorization header", http.StatusUnauthorized)
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			response.Send("", "Could not find bearer token in Authorization header", http.StatusUnauthorized)
			return
		}
		tokenRepo, _ := utility.GetContext("tokenRepo", r).(*repositories.TokenRepository)
		userID, err := tokenRepo.GetUserIDByToken(tokenString)
		if err != nil {
			response.Send("", err.Error(), err.HttpStatus())
			return
		}

		ctx := context.WithValue(r.Context(), "userName", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
