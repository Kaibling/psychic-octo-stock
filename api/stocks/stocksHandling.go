package stocks

import (
	"encoding/json"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/go-chi/chi"
)

func stockPost(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(w, r)
	userID := chi.URLParam(r, "userid")
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	_, err := userRepo.GetByID(userID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	var newStock models.Stock
	erra := json.NewDecoder(r.Body).Decode(&newStock)
	if erra != nil {
		response.Send("", err.Error(), http.StatusUnprocessableEntity)
		return

	}
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)
	if err := stockRepo.Add(&newStock); err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	stockRepo.AddStockToUser(newStock.ID, userID, newStock.Quantity)

	response.Send(newStock, "", http.StatusCreated)
}
func stocksGet(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(w, r)
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)
	userList, err := stockRepo.GetAll()
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(userList, "", http.StatusOK)
}
func stockPut(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(w, r)
	stockID := chi.URLParam(r, "id")
	var updateStock map[string]interface{}
	erra := json.NewDecoder(r.Body).Decode(&updateStock)
	if erra != nil {
		response.Send("", "post data not parsable", http.StatusUnprocessableEntity)
		return
	}

	updateStock["ID"] = stockID
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)

	stockRepo.UpdateWithMap(updateStock)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(loadedUser, "", http.StatusOK)
}
func stockDelete(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(w, r)
	stockID := chi.URLParam(r, "id")
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	if err := stockRepo.DeleteByObject(loadedUser); err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send("", "", http.StatusNoContent)

}
func stockGet(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(w, r)
	stockID := chi.URLParam(r, "id")
	stockRepo := utility.GetContext("stockRepo", r).(*repositories.StockRepository)
	loadedUser, err := stockRepo.GetByID(stockID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(loadedUser, "", http.StatusOK)
}
