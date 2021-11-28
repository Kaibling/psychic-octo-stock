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
