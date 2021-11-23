package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/dgrijalva/jwt-go"
)

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hmacSampleSecret, _ := utility.GetContext("hmacSecret", r).([]byte)
		auth := r.Header.Get("Authorization")
		if auth == "" {
			utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "Could not find Authorization header"}, http.StatusBadRequest)
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		if tokenString == auth {
			utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "Could not find bearer token in Authorization header"}, http.StatusBadRequest)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSampleSecret, nil
		})
		if err != nil {
			//log.Infof("Bad Thing happened: %s", err.Error())
			utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, http.StatusBadGateway)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			//log.Infoln("token invalid")
			utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "token invalid"}, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userName", claims["name"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
