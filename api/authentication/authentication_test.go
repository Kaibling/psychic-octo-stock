package authentication_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Kaibling/psychic-octo-stock/api"
	"github.com/Kaibling/psychic-octo-stock/api/authentication"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

var URL = "/api/v1/login"

func TestLogin(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "asd.asd@asd.as"}

	testLogin := authentication.UserLogin{
		Username: testUser.Username,
		Password: testUser.Password,
	}
	userRepo.Add(testUser)
	byte_User, _ := json.Marshal(testLogin)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	value := createResponse["data"]
	reponseToken := value.(string)
	assert.IsType(t, reponseToken, "string") //todo can be done better tested

}

func TestLoginWrongPassword(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "asd.asd@asd.as"}

	testLogin := authentication.UserLogin{
		Username: testUser.Username,
		Password: "asds",
	}
	userRepo.Add(testUser)
	byte_User, _ := json.Marshal(testLogin)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLoginWrongUsername(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "asd.asd@asd.as"}

	testLogin := authentication.UserLogin{
		Username: "asasada",
		Password: testUser.Password,
	}
	userRepo.Add(testUser)
	byte_User, _ := json.Marshal(testLogin)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}
