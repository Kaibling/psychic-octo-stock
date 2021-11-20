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
