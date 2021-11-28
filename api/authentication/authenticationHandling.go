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
	response := transmission.GetResponse(r)

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
	user, _ := userRepo.GetyName(userLogin.Username)

	tokenRepo, _ := utility.GetContext("tokenRepo", r).(*repositories.TokenRepository)
	token, err := tokenRepo.GenerateAndAddToken(user.ID, 0) //todo set valid date to 30 days or something
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(token, "", http.StatusOK)

}
