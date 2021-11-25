package stocks_test

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

var URL = "/api/v1/stocks"

func TestCreate(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	testUser := &models.User{Username: "Jack", Password: "abc123", Email: "abc.abc@abc.ab"}
	userRepo.Add(testUser)
	testStock := models.Stock{
		Name:     "Test",
		Quantity: 1223,
	}
	byteStock, _ := json.Marshal(testStock)
	w := performTestRequest(r, "POST", URL+"/users/"+testUser.ID, byteStock)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseObject, ok := value.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testStock.Name, reponseObject["name"])

}
func TestCreateMissingUser(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	testStock := models.Stock{
		Name: "Test",
	}
	byteStock, _ := json.Marshal(testStock)
	w := performTestRequest(r, "POST", URL+"/users/asda", byteStock)
	assert.Equal(t, http.StatusNotFound, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["message"]
	assert.True(t, exists)
	reponseObject, ok := value.(string)
	assert.True(t, ok)
	assert.Equal(t, reponseObject, "record not found")

}

func TestCreateNotUniqe(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123", Email: "abc.abc@abc.ab"}
	userRepo := repos["userRepo"].(*repositories.UserRepository)
	userRepo.Add(testUser)
	testStock := models.Stock{
		Name: "Test2",
	}
	byteStock, _ := json.Marshal(testStock)
	w := performTestRequest(r, "POST", URL+"/users/"+testUser.ID, byteStock)
	assert.Equal(t, http.StatusCreated, w.Code)
	//reapply for unique constrains violation
	w = performTestRequest(r, "POST", URL+"/users/"+testUser.ID, byteStock)
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

func TestUpdate(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	testStock := &models.Stock{Name: "Test3"}
	stockRepo := repos["stockRepo"].(*repositories.StockRepository)
	stockRepo.Add(testStock)
	stockID := testStock.ID

	updateObject := models.Stock{
		Name: "somethingNew",
	}
	updateBytes, _ := json.Marshal(updateObject)
	w := performTestRequest(r, "PUT", URL+"/"+stockID, updateBytes)
	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &updateResponse)
	assert.Nil(t, err)
	value2, exists := updateResponse["data"]
	assert.True(t, exists)
	reponseupdatedObject, ok := value2.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, updateObject.Name, reponseupdatedObject["name"])

}
func TestUpdateNoneExisting(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	objectID := "thisdoesnotexists"
	updateObject := models.Stock{
		Name: "somethingNew",
	}
	updateByts, _ := json.Marshal(updateObject)
	w := performTestRequest(r, "PUT", URL+objectID, updateByts)
	assert.Equal(t, http.StatusNotFound, w.Code)

}
func TestGetAll(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	stockRepo := repos["stockRepo"].(*repositories.StockRepository)
	testObject := &models.Stock{Name: "Test3"}
	stockRepo.Add(testObject)
	testObject2 := &models.Stock{Name: "Test4"}
	stockRepo.Add(testObject2)

	w := performTestRequest(r, "GET", URL, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseObjects, ok := value.([]interface{})
	assert.True(t, ok)

	object1 := reponseObjects[0].(map[string]interface{})
	assert.Equal(t, object1["name"], testObject.Name)

	object2 := reponseObjects[1].(map[string]interface{})
	assert.Equal(t, object2["name"], testObject2.Name)

}

func TestDelete(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	stockRepo := repos["stockRepo"].(*repositories.StockRepository)
	testObject := &models.Stock{Name: "Test3"}
	stockRepo.Add(testObject)
	objectID := testObject.ID

	deleteResponse := performTestRequest(r, "DELETE", URL+"/"+objectID, nil)
	assert.Equal(t, http.StatusNoContent, deleteResponse.Code)

	deleteResponse = performTestRequest(r, "DELETE", URL+"/"+objectID, nil)
	assert.Equal(t, http.StatusNotFound, deleteResponse.Code)
}

func TestDeleteNoneExisting(t *testing.T) {
	r, _, performTestRequest := api.TestAssemblyRoute()
	deleteResponse := performTestRequest(r, "DELETE", URL+"/adawfeefsse", nil)
	assert.Equal(t, http.StatusNotFound, deleteResponse.Code)

}

func TestGet(t *testing.T) {
	r, repos, performTestRequest := api.TestAssemblyRoute()
	stockRepo := repos["stockRepo"].(*repositories.StockRepository)
	testObject := &models.Stock{Name: "Test3"}
	stockRepo.Add(testObject)
	objectID := testObject.ID

	getResponse := performTestRequest(r, "GET", URL+"/"+objectID, nil)
	assert.Equal(t, http.StatusOK, getResponse.Code)

	var response map[string]interface{}
	err := json.Unmarshal(getResponse.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	//data empty
	object := value.(map[string]interface{})
	assert.Equal(t, object["name"], testObject.Name)

}
