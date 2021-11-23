package users

import (
	"encoding/json"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/go-chi/chi"
)

func userPost(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	erra := json.NewDecoder(r.Body).Decode(&newUser)
	if erra != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "post data not parsable"}, http.StatusUnprocessableEntity)
		return
	}
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	if err := userRepo.Add(&newUser); err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	//todo proper return schema
	newUser.Password = ""
	utility.SendResponse(w, r, &models.Envelope{Data: newUser, Message: ""}, http.StatusCreated)
	return
}
func usersGet(w http.ResponseWriter, r *http.Request) {
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)

	userList, err := userRepo.GetAll()
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	utility.SendResponse(w, r, &models.Envelope{Data: userList, Message: ""}, http.StatusOK)
	return
}
func userPut(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var updateUser map[string]interface{}
	erra := json.NewDecoder(r.Body).Decode(&updateUser)
	if erra != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: "post data not parsable"}, http.StatusUnprocessableEntity)
		return

	}

	updateUser["ID"] = userID
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)

	userRepo.UpdateWithMap(updateUser)
	loadedUser, err := userRepo.GetByID(userID)
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: ""}, err.HttpStatus())
		return
	}

	utility.SendResponse(w, r, &models.Envelope{Data: loadedUser, Message: ""}, http.StatusOK)
	return
}
func userDelete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	loadedUser, err := userRepo.GetByID(userID)
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: ""}, err.HttpStatus())
		return
	}
	if err := userRepo.DeleteByObject(loadedUser); err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	utility.SendResponse(w, r, nil, http.StatusNoContent)
	return
}
func userGet(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	userRepo := utility.GetContext("userRepo", r).(*repositories.UserRepository)
	loadedUser, err := userRepo.GetByID(userID)
	if err != nil {
		utility.SendResponse(w, r, &models.Envelope{Data: "", Message: err.Error()}, err.HttpStatus())
		return
	}
	utility.SendResponse(w, r, &models.Envelope{Data: loadedUser, Message: ""}, http.StatusOK)
	return
}
