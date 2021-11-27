package apimiddleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/dgrijalva/jwt-go"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := transmission.GetResponse(r)
		hmacSampleSecret, _ := utility.GetContext("hmacSecret", r).([]byte)
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

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSampleSecret, nil
		})
		if err != nil {
			response.Send("", err.Error(), http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response.Send("", "token invalid", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(r.Context(), "userName", claims["name"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
