package users_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kaibling/psychic-octo-stock/api"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string, jsonStr []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateUser(t *testing.T) {
	r := api.AssembleServer()
	testUser := models.User{
		Username: "Test",
		Email:    "abc@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performRequest(r, "POST", "/v1/users", byte_User)
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
	r := api.AssembleServer()
	testUser := models.User{
		Username: "Test2",
		Email:    "abc2@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performRequest(r, "POST", "/v1/users", byte_User)
	assert.Equal(t, http.StatusCreated, w.Code)
	//reapply for unique constrains violation
	w = performRequest(r, "POST", "/v1/users", byte_User)
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
	r := api.AssembleServer()
	testUser := models.User{
		Username: "Test3",
		Email:    "abc3@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performRequest(r, "POST", "/v1/users", byte_User)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseUser, ok := value.(map[string]interface{})
	assert.True(t, ok)
	userID := reponseUser["ID"].(string)

	updateUser := models.User{
		Address: "somethingNew",
	}
	updateByteUser, _ := json.Marshal(updateUser)
	w = performRequest(r, "PUT", "/v1/users/"+userID, updateByteUser)
	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
	assert.Nil(t, err)
	value2, exists := updateResponse["data"]
	assert.True(t, exists)
	reponseUser, ok = value2.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, updateUser.Address, reponseUser["address"])

}

func TestGetAllUser(t *testing.T) {
	r := api.AssembleServer()
	testUser := models.User{
		Username: "Test3",
		Email:    "abc3@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performRequest(r, "POST", "/v1/users", byte_User)
	assert.Equal(t, http.StatusCreated, w.Code)

	testUser2 := models.User{
		Username: "Test4",
		Email:    "abc3@abca.ac",
		Password: "abc123",
	}
	byte_User2, _ := json.Marshal(testUser2)
	w = performRequest(r, "POST", "/v1/users", byte_User2)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = performRequest(r, "GET", "/v1/users", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseUsers, ok := value.([]interface{})
	assert.True(t, ok)

	user1 := reponseUsers[0].(map[string]interface{})
	assert.Equal(t, user1["username"], testUser.Username)
	assert.Equal(t, user1["email"], testUser.Email)

	user2 := reponseUsers[1].(map[string]interface{})
	assert.Equal(t, user2["username"], testUser2.Username)
	assert.Equal(t, user2["email"], testUser2.Email)

	//userID := reponseUser["ID"].(string)

}

func TestDeleteUser(t *testing.T) {
	r := api.AssembleServer()
	testUser := models.User{
		Username: "Test3",
		Email:    "abc3@abc.ac",
		Password: "abc123",
	}
	byte_User, _ := json.Marshal(testUser)
	w := performRequest(r, "POST", "/v1/users", byte_User)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	value := response["data"]
	reponseUser := value.(map[string]interface{})
	userID := reponseUser["ID"].(string)

	deleteResponse := performRequest(r, "DELETE", "/v1/users/"+userID, nil)
	assert.Equal(t, http.StatusNoContent, deleteResponse.Code)

	deleteResponse = performRequest(r, "DELETE", "/v1/users/"+userID, nil)
	assert.Equal(t, http.StatusNotFound, deleteResponse.Code)
}

func TestDeleteNoneExistingUser(t *testing.T) {
	r := api.AssembleServer()
	deleteResponse := performRequest(r, "DELETE", "/v1/users/adawfeefsse", nil)
	assert.Equal(t, http.StatusNotFound, deleteResponse.Code)

}
