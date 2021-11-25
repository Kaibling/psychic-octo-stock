package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/repositories"
)

type UserLogin struct {
	Username string
	Password string
}

func login(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetOrCreateResponse(w, r)
	hmacSampleSecret := utility.GetContext("hmacSecret", r).([]byte)

	var userLogin UserLogin
	erra := json.NewDecoder(r.Body).Decode(&userLogin)
	if erra != nil {
		response.Send("", "post data not parsable", http.StatusUnprocessableEntity)
		return
	}
	userRepo, _ := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	psHash, err := userRepo.GetPWByName(userLogin.Username)
	if err != nil {
		response.Send("", "username/password incorrect", http.StatusUnauthorized)
		return
	}
	if !utility.ComparePasswords(psHash, userLogin.Password) {
		response.Send("", "username/password incorrect", http.StatusUnauthorized)
		return
	}

	token, erro := utility.GenerateToken(userLogin.Username, hmacSampleSecret)
	if erro != nil {
		response.Send("", err.Error(), http.StatusBadGateway)
		return
	}
	response.Send(token, "", http.StatusOK)

}
