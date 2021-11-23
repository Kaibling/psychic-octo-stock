package utility

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plain string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
func ComparePasswords(hashedPw string, comparePw string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPw), []byte(comparePw)); err != nil {
		log.Println("password compare error: " + err.Error())
		return false
	}
	return true
}

func BeautifyJson(data interface{}) string {
	b, _ := json.MarshalIndent(data, "", "  ")
	return string(b)
}

func GenerateToken(username string, hmacSampleSecret interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": username,
		"nbf":  time.Now().Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetContext(key string, r *http.Request) interface{} {
	parameter := r.Context().Value(key)
	if parameter == nil {
		panic(apierrors.NewClientError(errors.New("context parameter '" + key + "' missing")))
		//return "", apierrors.NewClientError(errors.New("context parameter '" + key + "' missing"))
	}
	return parameter
}

func SendResponse(w http.ResponseWriter, r *http.Request, data *models.Envelope, httpStatusCode int) {
	render.Status(r, httpStatusCode)
	render.Respond(w, r, data)
}
