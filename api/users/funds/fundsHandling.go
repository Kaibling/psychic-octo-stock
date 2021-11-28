package funds

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/apierrors"
	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/go-chi/chi"
)

type fundsData struct {
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Provider  string  `json:"provider"`
	Operation string  `json:"operation"` //CHARGE, WITHDRAW
}

func fundsPost(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	var newfundsData fundsData
	erra := json.NewDecoder(r.Body).Decode(&newfundsData)
	if erra != nil {
		response.Send("", "post data not parsable", http.StatusUnprocessableEntity)
		return
	}
	userID := chi.URLParam(r, "id")
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)

	err := executeFundOperation(userID, newfundsData, userRepo)
	if erra != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	loadedUser, _ := userRepo.GetByID(userID)
	response.Send(loadedUser, "", http.StatusOK)
}

func executeFundOperation(userID string, data fundsData, userRepo *repositories.UserRepository) apierrors.ApiError {
	switch operation := data.Operation; operation {
	case "CHARGE":
		//todo dont ignore Provider
		addFunds(userID, models.MonetaryUnit{Amount: data.Amount, Currency: data.Currency}, userRepo)

	case "WITHDRAW":
		return apierrors.NewGeneralError(errors.New("lol, nope"))
	default:
		return apierrors.NewGeneralError(errors.New("operation '" + operation + "' not supported"))

	}
	return nil
}

func addFunds(userID string, mu models.MonetaryUnit, userRepo *repositories.UserRepository) apierrors.ApiError {
	userRepo.AddFunds(userID, mu)
	return nil

}
