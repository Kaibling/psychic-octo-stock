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

func performRequest(r http.Handler, method, path string, jsonStr []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonStr))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
