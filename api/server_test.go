package api_test

import (
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
	w := performTestRequest(r, "GET", URL+"a", nil)
	assert.Equal(t, http.StatusNotFound, w.Code)

}
