package users_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Kaibling/psychic-octo-stock/api"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

var URL = "/api/v1/users"

func TestCreateUser(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	testUser := models.User{
		Username: "Test",
		Email:    "abc@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseUser, ok := value.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testUser.Username, reponseUser["username"])
	assert.Equal(t, testUser.Email, reponseUser["email"])

}
func TestCreateUserNotUniqe(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	testUser := models.User{
		Username: "Test2",
		Email:    "abc2@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	//reapply for unique constrains violation
	w = performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	//data empty
	assert.Equal(t, "", value)

	//something in the message
	message, exists := response["message"]
	assert.True(t, exists)
	_, ok := message.(string)
	assert.True(t, ok)
	//todo maybe compare error message

}

func TestUpdateUser(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	testUser := models.User{
		Username: "Test3",
		Email:    "abc3@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseUser, ok := value.(map[string]interface{})
	assert.True(t, ok)
	userID := reponseUser["ID"].(string)

	//updateUser := models.User{
	//	Address: "somethingNew",
	//}
	updateUser := map[string]interface{}{
		"Address": "somethingNew",
	}
	updateByteUser, _ := json.Marshal(updateUser)
	w = performTestRequest(r, "PUT", URL+"/"+userID, updateByteUser, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
	assert.Nil(t, err)
	value2, exists := updateResponse["data"]
	assert.True(t, exists)
	reponseUser, ok = value2.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, updateUser["Address"], reponseUser["address"])

}
func TestUpdateNoneExistingUser(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	userID := "thisdoesnotexists"
	updateUser := models.User{
		Address: "somethingNew",
	}
	updateByteUser, _ := json.Marshal(updateUser)
	w := performTestRequest(r, "PUT", URL+"/"+userID, updateByteUser, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)

}
func TestGetAllUser(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	testUser := models.User{
		Username: "Test3",
		Email:    "abc3@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	testUser2 := models.User{
		Username: "Test4",
		Email:    "abc3@abca.ac",
		Password: "abc123",
	}
	byte_User2, _ := json.Marshal(testUser2)
	w = performTestRequest(r, "POST", URL, byte_User2, nil)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = performTestRequest(r, "GET", URL, nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseUsers, ok := value.([]interface{})
	assert.True(t, ok)

	user1 := reponseUsers[1].(map[string]interface{})
	assert.Equal(t, user1["username"], testUser.Username)
	assert.Equal(t, user1["email"], testUser.Email)

	user2 := reponseUsers[2].(map[string]interface{})
	assert.Equal(t, user2["username"], testUser2.Username)
	assert.Equal(t, user2["email"], testUser2.Email)

	//userID := reponseUser["ID"].(string)

}

func TestDeleteUser(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	testUser := models.User{
		Username: "Test3",
		Email:    "abc3@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	value := response["data"]
	reponseUser := value.(map[string]interface{})
	userID := reponseUser["ID"].(string)

	deleteResponse := performTestRequest(r, "DELETE", URL+"/"+userID, nil, nil)
	assert.Equal(t, http.StatusNoContent, deleteResponse.Code)

	deleteResponse = performTestRequest(r, "DELETE", URL+"/"+userID, nil, nil)
	assert.Equal(t, http.StatusNotFound, deleteResponse.Code)
}

func TestDeleteNoneExistingUser(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	deleteResponse := performTestRequest(r, "DELETE", URL+"/adawfeefsse", nil, nil)
	assert.Equal(t, http.StatusNotFound, deleteResponse.Code)

}

func TestGetUser(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	testUser := models.User{
		Username: "Test3",
		Email:    "abc3@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performTestRequest(r, "POST", URL, byte_User, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	value := createResponse["data"]
	reponseUser := value.(map[string]interface{})
	userID := reponseUser["ID"].(string)

	getResponse := performTestRequest(r, "GET", URL+"/"+userID, nil, nil)
	assert.Equal(t, http.StatusOK, getResponse.Code)

	var response map[string]interface{}
	err := json.Unmarshal(getResponse.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	//data empty
	user := value.(map[string]interface{})
	assert.Equal(t, user["email"], testUser.Email)
	assert.Equal(t, user["username"], testUser.Username)
}

func TestUserFunds(t *testing.T) {
	_, repos, _ := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "a@a.a"}
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	userRepo.Add(testUser)
	userFunds, err := userRepo.FundsByID(testUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, 1234.0, userFunds)
}
