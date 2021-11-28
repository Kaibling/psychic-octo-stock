package tokens

import (
	"encoding/json"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/go-chi/chi"
)

type tokenData struct {
	ValidUntil uint64 `json:"valid_until"`
}

func tokenPost(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	var newTokenData tokenData
	erra := json.NewDecoder(r.Body).Decode(&newTokenData)
	if erra != nil {
		response.Send("", "post data not parsable", http.StatusUnprocessableEntity)
		return
	}
	userID := chi.URLParam(r, "id")
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	if _, err := userRepo.GetByID(userID); err != nil {
		response.Send("", "user not found", http.StatusNotFound)
		return
	}
	tokenRepo := utility.GetContext("tokenRepo", r).(*repositories.TokenRepository)
	tokenString, err := tokenRepo.GenerateAndAddToken(userID, int64(newTokenData.ValidUntil))
	if erra != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(tokenString, "", http.StatusCreated)
}

func tokenPut(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	//userID := chi.URLParam(r, "id")
	tokenID := chi.URLParam(r, "tid")
	var updateToken map[string]interface{}
	erra := json.NewDecoder(r.Body).Decode(&updateToken)
	if erra != nil {
		response.Send("", "post data not parsable", http.StatusUnprocessableEntity)
		return
	}

	updateToken["ID"] = tokenID
	tokenRepo := utility.GetContext("tokenRepo", r).(*repositories.TokenRepository)

	tokenRepo.UpdateWithMap(updateToken)
	loadedToken, err := tokenRepo.GetByID(tokenID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(loadedToken, "", http.StatusOK)
}
func tokenDelete(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	//userID := chi.URLParam(r, "id")
	tokenID := chi.URLParam(r, "tid")
	tokenRepo := utility.GetContext("tokenRepo", r).(*repositories.TokenRepository)
	if err := tokenRepo.DeleteByID(tokenID); err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send("", "", http.StatusNoContent)

}

func tokensGet(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	userID := chi.URLParam(r, "id")
	tokenRepo := utility.GetContext("tokenRepo", r).(*repositories.TokenRepository)
	loadedToken, err := tokenRepo.GetAll(userID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(loadedToken, "", http.StatusOK)
}
