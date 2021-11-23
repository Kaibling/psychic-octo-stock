package stocks

import (
	"encoding/json"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/go-chi/chi"
)

func stockPost(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userid")
	// if err != nil {
	// 	c.JSON(err.HttpStatus(), models.Envelope{Data: "", Message: err.Error()})
	// 	return
	// }
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	_, err := userRepo.GetByID(userID)
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	var newStock models.Stock
	erra := json.NewDecoder(r.Body).Decode(&newStock)
	if erra != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: erra.Error()}, http.StatusUnprocessableEntity)
		return

	}
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)
	if err := stockRepo.Add(&newStock); err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	//add ownership to an user
	stockRepo.AddStockToUser(newStock.ID, userID, newStock.Quantity)

	env := models.Envelope{Data: newStock, Message: ""}
	utility.SendResponse(w, r, &env, http.StatusCreated)
}
func stocksGet(w http.ResponseWriter, r *http.Request) {
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)
	userList, err := stockRepo.GetAll()
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	env := models.Envelope{Data: userList, Message: ""}
	utility.SendResponse(w, r, &env, http.StatusOK)
}
func stockPut(w http.ResponseWriter, r *http.Request) {
	stockID := chi.URLParam(r, "id")
	var updateStock map[string]interface{}
	erra := json.NewDecoder(r.Body).Decode(&updateStock)
	if erra != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "post data not parsable"}, http.StatusUnprocessableEntity)
		return
	}

	updateStock["ID"] = stockID
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)

	stockRepo.UpdateWithMap(updateStock)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	env := models.Envelope{Data: loadedUser, Message: ""}
	utility.SendResponse(w, r, &env, http.StatusOK)
}
func stockDelete(w http.ResponseWriter, r *http.Request) {
	stockID := chi.URLParam(r, "id")
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	if err := stockRepo.DeleteByObject(loadedUser); err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	utility.SendResponse(w, r, nil, http.StatusNoContent)
}
func stockGet(w http.ResponseWriter, r *http.Request) {
	stockID := chi.URLParam(r, "id")
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	env := models.Envelope{Data: loadedUser, Message: ""}
	utility.SendResponse(w, r, &env, http.StatusOK)
}
