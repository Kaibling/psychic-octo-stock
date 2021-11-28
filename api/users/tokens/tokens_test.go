package tokens_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Kaibling/psychic-octo-stock/api"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/stretchr/testify/assert"
)

// func TestMain(m *testing.M) {
// 	log.SetOutput(ioutil.Discard)
// 	os.Exit(m.Run())
// }

var URL = "/api/v1/users"

func TestGenerateToken(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "asd.asd@asd.as"}
	userRepo.Add(testUser)

	requestdate := map[string]interface{}{
		"valid_until": 123456,
	}

	byteData, _ := json.Marshal(requestdate)
	w := performTestRequest(r, "POST", URL+"/"+testUser.ID+"/tokens", byteData, nil)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	value := createResponse["data"]
	reponseToken := value.(string)
	assert.IsType(t, reponseToken, "string")

}

func TestGenerateTokenWrongUser(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "asd.asd@asd.as"}
	userRepo.Add(testUser)

	requestdate := map[string]interface{}{
		"valid_until": 123456,
	}
	byteData, _ := json.Marshal(requestdate)
	w := performTestRequest(r, "POST", URL+"/ss/tokens", byteData, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllTokens(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "asd.asd@asd.as"}
	userRepo.Add(testUser)

	tokenRepo := repos["tokenRepo"].(*repositories.TokenRepository)

	token1, err := tokenRepo.GenerateAndAddToken(testUser.ID, 0)
	assert.Nil(t, err)
	token2, err := tokenRepo.GenerateAndAddToken(testUser.ID, 123456)
	assert.Nil(t, err)

	w := performTestRequest(r, "GET", URL+"/"+testUser.ID+"/tokens", nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	value := createResponse["data"]
	reponseObjects, ok := value.([]interface{})
	assert.True(t, ok)

	object1 := reponseObjects[0].(map[string]interface{})
	assert.Equal(t, object1["user_id"], testUser.ID)
	assert.Equal(t, object1["token"], token1)
	assert.Equal(t, float64(0), object1["valid_until"])

	object2 := reponseObjects[1].(map[string]interface{})
	assert.Equal(t, object2["user_id"], testUser.ID)
	assert.Equal(t, object2["token"], token2)
	assert.Equal(t, float64(123456), object2["valid_until"])

}

func TestDeleteToken(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "asd.asd@asd.as"}
	userRepo.Add(testUser)

	tokenRepo := repos["tokenRepo"].(*repositories.TokenRepository)
	_, err := tokenRepo.GenerateAndAddToken(testUser.ID, 123456)
	assert.Nil(t, err)

	w := performTestRequest(r, "GET", URL+"/"+testUser.ID+"/tokens", nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	var Response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &Response)
	value := Response["data"]
	reponseObjects, ok := value.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, len(reponseObjects), 1)

	object1 := reponseObjects[0].(map[string]interface{})
	assert.Equal(t, object1["user_id"], testUser.ID)
	tokenID := object1["ID"].(string)

	w = performTestRequest(r, "DELETE", URL+"/"+testUser.ID+"/tokens/"+tokenID, nil, nil)
	assert.Equal(t, http.StatusNoContent, w.Code)

	w = performTestRequest(r, "GET", URL+"/"+testUser.ID+"/tokens", nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	json.Unmarshal(w.Body.Bytes(), &Response)
	value = Response["data"]
	reponseObjects, ok = value.([]interface{})
	assert.True(t, ok)
	assert.Equal(t, len(reponseObjects), 0)
}

func TestPutToken(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 1234, Email: "asd.asd@asd.as"}
	userRepo.Add(testUser)

	tokenRepo := repos["tokenRepo"].(*repositories.TokenRepository)

	token, err := tokenRepo.GenerateAndAddToken(testUser.ID, 123456)
	assert.Nil(t, err)

	w := performTestRequest(r, "GET", URL+"/"+testUser.ID+"/tokens", nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	var Response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &Response)
	value := Response["data"]
	reponseObjects, ok := value.([]interface{})
	assert.True(t, ok)

	responseToken1 := reponseObjects[0].(map[string]interface{})
	assert.Equal(t, len(reponseObjects), 1)
	assert.Equal(t, token, responseToken1["token"].(string))
	assert.Equal(t, testUser.ID, responseToken1["user_id"].(string))
	tokenID := responseToken1["ID"].(string)

	updateToken := map[string]interface{}{
		"comment": "testcomment",
	}

	byteData, _ := json.Marshal(updateToken)
	w = performTestRequest(r, "PUT", URL+"/"+testUser.ID+"/tokens/"+tokenID, byteData, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updateResponse)
	assert.Nil(t, err)
	value, exists := updateResponse["data"]
	assert.True(t, exists)
	reponseupdatedObject, ok := value.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "testcomment", reponseupdatedObject["comment"])

}
