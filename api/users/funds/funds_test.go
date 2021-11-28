package funds_test

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

func TestChargeFund(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 200, Email: "asd.asd@asd.as", Currency: "EUR"}
	userRepo.Add(testUser)

	requestFunds := map[string]interface{}{
		"Amount":    123,
		"Currency":  "EUR",
		"Provider":  "",
		"Operation": "CHARGE",
	}

	byteData, _ := json.Marshal(requestFunds)
	w := performTestRequest(r, "POST", URL+"/"+testUser.ID+"/funds", byteData, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	value := createResponse["data"]
	responseUser := value.(map[string]interface{})
	assert.Equal(t, float64(323), responseUser["funds"])
	assert.Equal(t, testUser.Currency, responseUser["currency"])
}

func TestChargeFundDifferentCUrrency(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Funds: 200, Email: "asd.asd@asd.as", Currency: "EUR"}
	userRepo.Add(testUser)

	requestFunds := map[string]interface{}{
		"Amount":    123,
		"Currency":  "AED",
		"Provider":  "",
		"Operation": "CHARGE",
	}

	byteData, _ := json.Marshal(requestFunds)
	w := performTestRequest(r, "POST", URL+"/"+testUser.ID+"/funds", byteData, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	value := createResponse["data"]
	responseUser := value.(map[string]interface{})
	assert.Equal(t, float64(229.1018), responseUser["funds"])
	assert.Equal(t, testUser.Currency, responseUser["currency"])
}
