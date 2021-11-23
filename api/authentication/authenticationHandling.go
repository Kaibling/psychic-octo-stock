package authentication

import (
	"encoding/json"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
)

type UserLogin struct {
	Username string
	Password string
}

func login(w http.ResponseWriter, r *http.Request) {

	hmacSampleSecret := utility.GetContext("hmacSecret", r).([]byte)

	var userLogin UserLogin
	erra := json.NewDecoder(r.Body).Decode(&userLogin)
	if erra != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "post data not parsable"}, http.StatusUnprocessableEntity)
		return
	}
	userRepo, _ := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	psHash, err := userRepo.GetPWByName(userLogin.Username)
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "username/password incorrect"}, http.StatusUnauthorized)
		return
	}
	if !utility.ComparePasswords(psHash, userLogin.Password) {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "username/password incorrect"}, http.StatusUnauthorized)
		return
	}

	token, erro := utility.GenerateToken(userLogin.Username, hmacSampleSecret)
	if erro != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, http.StatusBadGateway)
		return
	}
	utility.SendResponse(w, r, &models.Envelope{Data: token, Message: ""}, http.StatusOK)
}
