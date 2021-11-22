package middleware

import (
	"fmt"
	"strings"

	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authorization(c *gin.Context) {

	hmacSampleSecret := c.MustGet("hmacSecret").([]byte)
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		c.JSON(403, models.Envelope{Data: "", Message: "Could not find Authorization header"})
		c.Abort()
		return

	}
	tokenString := strings.TrimPrefix(auth, "Bearer ")
	if tokenString == auth {
		c.JSON(403, models.Envelope{Data: "", Message: "Could not find bearer token in Authorization header"})
		c.Abort()
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
		c.JSON(500, models.Envelope{Data: "failed", Message: err.Error()})
		c.Abort()
		return

	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		//log.Infoln("token invalid")
		c.JSON(403, models.Envelope{Data: "failed", Message: "token invalid"})
		c.Abort()
		return
	}
	c.Set("userName", claims["name"])
	c.Next()
}
