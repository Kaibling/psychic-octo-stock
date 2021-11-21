package transactions_test

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

var URL = "/v1/transactions"

func TestCreate(t *testing.T) {
	r, userRepo, stockRepo, _ := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123"}
	userRepo.Add(testUser)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	testObject := models.Transaction{
		UserID:   testUser.ID,
		StockID:  testStock.ID,
		Quantity: 12,
		Type:     "SELL",
	}
	byteObject, _ := json.Marshal(testObject)
	w := performRequest(r, "POST", URL, byteObject)
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseObject, ok := value.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, testObject.Type, reponseObject["type"])
	assert.Equal(t, testObject.UserID, reponseObject["userID"])
	assert.Equal(t, testObject.StockID, reponseObject["stockID"])
	//todo fix
	//assert.Equal(t, testObject.Quantity, reponseObject["quantity"])

}

// func TestCreateNotenoughStocks(t *testing.T) {
// 	r, userRepo, stockRepo, _ := api.TestAssemblyRoute()
// 	testUser := &models.User{Username: "Jack", Password: "abc123"}
// 	userRepo.Add(testUser)
// 	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
// 	stockRepo.Add(testStock)
// 	testObject := models.Transaction{
// 		UserID:   testUser.ID,
// 		StockID:  testStock.ID,
// 		Quantity: 12,
// 		Type:     "SELL",
// 	}
// 	byteObject, _ := json.Marshal(testObject)
// 	w := performRequest(r, "POST", URL, byteObject)
// 	assert.Equal(t, http.StatusCreated, w.Code)
// 	//reapply for unique constrains violation
// 	w = performRequest(r, "POST", URL, byteObject)
// 	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

// 	var response map[string]interface{}
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.Nil(t, err)
// 	value, exists := response["data"]
// 	assert.True(t, exists)
// 	//data empty
// 	assert.Equal(t, "", value)

// 	//something in the message
// 	message, exists := response["message"]
// 	assert.True(t, exists)
// 	_, ok := message.(string)
// 	assert.True(t, ok)
// 	//todo maybe compare error message

// }

func TestGetAll(t *testing.T) {
	r, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123"}
	userRepo.Add(testUser)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	testObject := models.Transaction{
		UserID:   testUser.ID,
		StockID:  testStock.ID,
		Quantity: 12,
		Type:     "SELL",
	}
	transactionRepo.Add(&testObject)
	testObject2 := models.Transaction{
		UserID:   testUser.ID,
		StockID:  testStock.ID,
		Quantity: 123,
		Type:     "BUY",
	}
	transactionRepo.Add(&testObject2)
	w := performRequest(r, "GET", URL, nil)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	reponseObjects, ok := value.([]interface{})
	assert.True(t, ok)

	object1 := reponseObjects[0].(map[string]interface{})
	assert.Equal(t, object1["userID"], testObject.UserID)
	assert.Equal(t, object1["stockID"], testObject.StockID)
	//assert.Equal(t, object1["quantity"], testObject.Quantity)
	assert.Equal(t, object1["type"], testObject.Type)

	object2 := reponseObjects[1].(map[string]interface{})
	assert.Equal(t, object2["userID"], testObject2.UserID)
	assert.Equal(t, object2["stockID"], testObject2.StockID)
	//assert.Equal(t, object2["quantity"], testObject2.Quantity)
	assert.Equal(t, object2["type"], testObject2.Type)

}

func TestDelete(t *testing.T) {
	r, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123"}
	userRepo.Add(testUser)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	testObject := models.Transaction{
		UserID:   testUser.ID,
		StockID:  testStock.ID,
		Quantity: 12,
		Type:     "SELL",
	}
	transactionRepo.Add(&testObject)
	objectID := testObject.ID

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
	r, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123"}
	userRepo.Add(testUser)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	testObject := models.Transaction{
		UserID:   testUser.ID,
		StockID:  testStock.ID,
		Quantity: 12,
		Type:     "SELL",
	}
	transactionRepo.Add(&testObject)
	objectID := testObject.ID

	getResponse := performRequest(r, "GET", URL+"/"+objectID, nil)
	assert.Equal(t, http.StatusOK, getResponse.Code)

	var response map[string]interface{}
	err := json.Unmarshal(getResponse.Body.Bytes(), &response)
	assert.Nil(t, err)
	value, exists := response["data"]
	assert.True(t, exists)
	//data empty
	object := value.(map[string]interface{})
	//assert.Equal(t, object["Quantity"], testObject.Quantity)
	assert.Equal(t, object["Type"], testObject.Type)

}
