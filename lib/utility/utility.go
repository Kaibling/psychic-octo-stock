package utility

import (
	"encoding/json"
	"log"

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