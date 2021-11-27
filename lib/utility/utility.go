package utility

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/config"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plain string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
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
	}
	return parameter
}

func CurrencyConverter(mu models.MonetaryUnit, currency string) models.MonetaryUnit {
	//todo things
	return models.MonetaryUnit{Amount: mu.Amount, Currency: currency}
}

func AddAndConvertFunds(mu1 models.MonetaryUnit, mu2 models.MonetaryUnit) models.MonetaryUnit {
	//if both have the same currency, just add it
	//if one has the default one, take the default
	//if nobody has default, convert both to default
	if mu1.Currency == mu2.Currency {
		return models.MonetaryUnit{Amount: mu1.Amount + mu2.Amount, Currency: mu1.Currency}
	} else if mu1.Currency == config.Config.Currency || mu2.Currency == config.Config.Currency {
		if mu1.Currency == config.Config.Currency {
			mu2Converted := CurrencyConverter(mu2, config.Config.Currency)
			return models.MonetaryUnit{Amount: mu1.Amount + mu2Converted.Amount, Currency: config.Config.Currency}
		}
		mu1Converted := CurrencyConverter(mu1, config.Config.Currency)
		return models.MonetaryUnit{Amount: mu2.Amount + mu1Converted.Amount, Currency: config.Config.Currency}
	}
	mu1Converted := CurrencyConverter(mu1, config.Config.Currency)
	mu2Converted := CurrencyConverter(mu2, config.Config.Currency)
	return models.MonetaryUnit{Amount: mu1Converted.Amount + mu2Converted.Amount, Currency: config.Config.Currency}
}

func SubtractAndConvertFunds(mu1 models.MonetaryUnit, mu2 models.MonetaryUnit) models.MonetaryUnit {
	//if both have the same currency, just add it
	//if one has the default one, take the default
	//if nobody has default, convert both to default
	if mu1.Currency == mu2.Currency {
		return models.MonetaryUnit{Amount: mu1.Amount - mu2.Amount, Currency: mu1.Currency}
	} else if mu1.Currency == config.Config.Currency || mu2.Currency == config.Config.Currency {
		if mu1.Currency == config.Config.Currency {
			mu2Converted := CurrencyConverter(mu2, config.Config.Currency)
			return models.MonetaryUnit{Amount: mu1.Amount - mu2Converted.Amount, Currency: config.Config.Currency}
		}
		mu1Converted := CurrencyConverter(mu1, config.Config.Currency)
		return models.MonetaryUnit{Amount: mu2.Amount - mu1Converted.Amount, Currency: config.Config.Currency}
	}
	mu1Converted := CurrencyConverter(mu1, config.Config.Currency)
	mu2Converted := CurrencyConverter(mu2, config.Config.Currency)
	return models.MonetaryUnit{Amount: mu1Converted.Amount - mu2Converted.Amount, Currency: config.Config.Currency}
}
