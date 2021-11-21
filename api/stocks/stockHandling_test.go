package stocks_test

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

var URL = "/v1/stocks"

func TestCreate(t *testing.T) {
	r := api.AssembleServer()
	testStock := models.Stock{
		Name: "Test",
	}
	byteStock, _ := json.Marshal(testStock)
	w := performRequest(r, "POST", URL, byteStock)
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

func TestCreateNotUniqe(t *testing.T) {
	r := api.AssembleServer()
	testStock := models.Stock{
		Name: "Test2",
	}
	byteStock, _ := json.Marshal(testStock)
	w := performRequest(r, "POST", URL, byteStock)
	assert.Equal(t, http.StatusCreated, w.Code)
	//reapply for unique constrains violation
	w = performRequest(r, "POST", URL, byteStock)
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
	r := api.AssembleServer()
	testStock := models.Stock{
		Name: "Test3",
	}
	byteStock, _ := json.Marshal(testStock)
	w := performRequest(r, "POST", URL, byteStock)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseObject, ok := value.(map[string]interface{})
	assert.True(t, ok)
	stockID := reponseObject["ID"].(string)

	updateObject := models.Stock{
		Name: "somethingNew",
	}
	updateBytes, _ := json.Marshal(updateObject)
	w = performRequest(r, "PUT", URL+"/"+stockID, updateBytes)
	assert.Equal(t, http.StatusOK, w.Code)

	var updateResponse map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &updateResponse)
	assert.Nil(t, err)
	value2, exists := updateResponse["data"]
	assert.True(t, exists)
	reponseupdatedObject, ok := value2.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, updateObject.Name, reponseupdatedObject["name"])

}
func TestUpdateNoneExisting(t *testing.T) {
	r := api.AssembleServer()
	objectID := "thisdoesnotexists"
	updateObject := models.Stock{
		Name: "somethingNew",
	}
	updateByts, _ := json.Marshal(updateObject)
	w := performRequest(r, "PUT", URL+objectID, updateByts)
	assert.Equal(t, http.StatusNotFound, w.Code)

}
func TestGetAll(t *testing.T) {
	r := api.AssembleServer()
	testObject := models.Stock{
		Name: "Test3",
	}
	byteObject, _ := json.Marshal(testObject)
	w := performRequest(r, "POST", URL, byteObject)
	assert.Equal(t, http.StatusCreated, w.Code)

	testObject2 := models.Stock{
		Name: "Test4",
	}
	byteObject2, _ := json.Marshal(testObject2)
	w = performRequest(r, "POST", URL, byteObject2)
	assert.Equal(t, http.StatusCreated, w.Code)

	w = performRequest(r, "GET", URL, nil)
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
	r := api.AssembleServer()
	testObject := models.Stock{
		Name: "Test3",
	}
	byteObject, _ := json.Marshal(testObject)
	w := performRequest(r, "POST", URL, byteObject)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	value := response["data"]
	reponseObject := value.(map[string]interface{})
	objectID := reponseObject["ID"].(string)

	deleteResponse := performRequest(r, "DELETE", URL+"/"+objectID, nil)
	assert.Equal(t, http.StatusNoContent, deleteResponse.Code)

	deleteResponse = performRequest(r, "DELETE", URL+"/"+objectID, nil)
	assert.Equal(t, http.StatusNotFound, deleteResponse.Code)
}

func TestDeleteNoneExisting(t *testing.T) {
	r := api.AssembleServer()
	deleteResponse := performRequest(r, "DELETE", URL+"/adawfeefsse", nil)
	assert.Equal(t, http.StatusNotFound, deleteResponse.Code)

}

func TestGet(t *testing.T) {
	r := api.AssembleServer()
	testObject := models.Stock{
		Name: "Test3",
	}
	byteObject, _ := json.Marshal(testObject)
	w := performRequest(r, "POST", URL, byteObject)
	assert.Equal(t, http.StatusCreated, w.Code)
	var createResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResponse)
	value := createResponse["data"]
	reponseObject := value.(map[string]interface{})
	objectID := reponseObject["ID"].(string)

	getResponse := performRequest(r, "GET", URL+"/"+objectID, nil)
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
