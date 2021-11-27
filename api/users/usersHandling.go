package users

import (
	"encoding/json"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/go-chi/chi"
)

func userPost(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	var newUser models.User
	erra := json.NewDecoder(r.Body).Decode(&newUser)
	if erra != nil {
		response.Send("", "post data not parsable", http.StatusUnprocessableEntity)
		return
	}
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	if err := userRepo.Add(&newUser); err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	//todo proper return schema
	newUser.Password = ""
	response.Send(newUser, "", http.StatusCreated)
}

func usersGet(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)

	userList, err := userRepo.GetAll()
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(userList, "", http.StatusOK)
}

func userPut(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	userID := chi.URLParam(r, "id")
	var updateUser map[string]interface{}
	erra := json.NewDecoder(r.Body).Decode(&updateUser)
	if erra != nil {
		response.Send("", "post data not parsable", http.StatusUnprocessableEntity)
		return

	}

	updateUser["ID"] = userID
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)

	userRepo.UpdateWithMap(updateUser)
	loadedUser, err := userRepo.GetByID(userID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}

	response.Send(loadedUser, "", http.StatusOK)
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	userID := chi.URLParam(r, "id")
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	loadedUser, err := userRepo.GetByID(userID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	if err := userRepo.DeleteByObject(loadedUser); err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send("", "", http.StatusNoContent)
}

func userGet(w http.ResponseWriter, r *http.Request) {
	response := transmission.GetResponse(r)
	userID := chi.URLParam(r, "id")
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	loadedUser, err := userRepo.GetByID(userID)
	if err != nil {
		response.Send("", err.Error(), err.HttpStatus())
		return
	}
	response.Send(loadedUser, "", http.StatusOK)
}
