package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
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

func GetContext(key string, r *http.Request) interface{} {
	parameter := r.Context().Value(key)
	if parameter == nil {
		panic(apierrors.NewClientError(errors.New("context parameter '" + key + "' missing")))
	}
	return parameter
}

func GetRequest(url string) []byte {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
