package api_test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/Kaibling/psychic-octo-stock/api"
	"github.com/stretchr/testify/assert"
)

var URL = "/api/v1/users"

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestNoneExistingRoute(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	w := performTestRequest(r, "GET", URL+"a", nil, nil)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestSendRequestID(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	expectedRequestID := "abopt34f"
	w := performTestRequest(r, "GET", URL, nil, &map[string]string{"X-REQUEST-ID": expectedRequestID})
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	value := response["request_id"].(string)
	assert.Equal(t, value, expectedRequestID)
}

func TestGeneratedRequestID(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	w := performTestRequest(r, "GET", URL, nil, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	value := response["request_id"].(string)
	assert.Greater(t, len(value), 4)
}
