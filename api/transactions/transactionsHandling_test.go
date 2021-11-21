package transactions_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kaibling/psychic-octo-stock/api"
	"github.com/Kaibling/psychic-octo-stock/api/transactions"
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
		SellerID: testUser.ID,
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
	assert.Equal(t, testObject.SellerID, reponseObject["sellerID"])
	assert.Equal(t, testObject.StockID, reponseObject["stockID"])
	assert.Equal(t, testObject.Quantity, int(reponseObject["quantity"].(float64)))

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
		SellerID: testUser.ID,
		StockID:  testStock.ID,
		Quantity: 12,
		Type:     "SELL",
	}
	transactionRepo.Add(&testObject)
	testObject2 := models.Transaction{
		BuyerID:  testUser.ID,
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
	assert.Equal(t, object1["buyerID"], testObject.BuyerID)
	assert.Equal(t, object1["stockID"], testObject.StockID)
	assert.Equal(t, int(object1["quantity"].(float64)), testObject.Quantity)
	assert.Equal(t, object1["type"], testObject.Type)

	object2 := reponseObjects[1].(map[string]interface{})
	assert.Equal(t, object2["sellerID"], testObject2.SellerID)

	assert.Equal(t, object2["stockID"], testObject2.StockID)
	assert.Equal(t, int(object2["quantity"].(float64)), testObject2.Quantity)
	assert.Equal(t, object2["type"], testObject2.Type)

}

func TestDelete(t *testing.T) {
	r, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123"}
	userRepo.Add(testUser)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	testObject := models.Transaction{
		SellerID: testUser.ID,
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
		SellerID: testUser.ID,
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
	object := value.(map[string]interface{})
	assert.Equal(t, int(object["quantity"].(float64)), testObject.Quantity)
	assert.Equal(t, object["type"], testObject.Type)

}

func TestAtomicFunction(t *testing.T) {
	_, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123", Address: "abc-street 123"}
	userRepo.Add(testUser)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	atomicExecutionArray := []interface{}{}
	updateTestUser := map[string]interface{}{"Address": "cba-street 321"}
	var updateTestUserQuery interface{} = "id = ?"
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.User{}, updateTestUser, updateTestUserQuery, []interface{}{testUser.ID}})
	updateTestStock := map[string]interface{}{"Quantity": 321}
	var updateTestStockQuery interface{} = "id = ?"
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.Stock{}, updateTestStock, updateTestStockQuery, []interface{}{testStock.ID}})
	err := transactionRepo.ExecuteTransaction(atomicExecutionArray)
	assert.Nil(t, err)
	checkUser, err := userRepo.GetByID(testUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, "cba-street 321", checkUser.Address)

	checkStock, err := stockRepo.GetByID(testStock.ID)
	assert.Nil(t, err)
	assert.Equal(t, 321, checkStock.Quantity)
}

func TestAtomicFunctionRollback(t *testing.T) {
	_, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testUser := &models.User{Username: "Jack", Password: "abc123", Address: "abc-street 123"}
	userRepo.Add(testUser)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	atomicExecutionArray := []interface{}{}
	updateTestUser := map[string]interface{}{"Address": "cba-street 321"}
	var updateTestUserQuery interface{} = "id = ?"
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.User{}, updateTestUser, updateTestUserQuery, []interface{}{testUser.ID}})
	updateTestStock := map[string]interface{}{"Quantity": 321}
	var updateTestStockQuery interface{} = "id = ?"
	atomicExecutionArray = append(atomicExecutionArray, []interface{}{models.Stock{}, updateTestStock, updateTestStockQuery, []interface{}{"noneexisting"}})
	err := transactionRepo.ExecuteTransaction(atomicExecutionArray)
	assert.NotNil(t, err)
	checkUser, err := userRepo.GetByID(testUser.ID)
	assert.Nil(t, err)
	assert.Equal(t, "abc-street 123", checkUser.Address)

	checkStock, err := stockRepo.GetByID(testStock.ID)
	assert.Nil(t, err)
	assert.Equal(t, 123, checkStock.Quantity)
}

// status

func TestStatusSetActive(t *testing.T) {

	_, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testSeller := &models.User{Username: "Jack", Password: "abc123"}
	userRepo.Add(testSeller)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	testObject := models.Transaction{
		SellerID: testSeller.ID,
		StockID:  testStock.ID,
		Quantity: 12,
		Type:     "SELL",
		Status:   "PENDING",
	}
	transactionRepo.Add(&testObject)
	assert.Equal(t, testObject.Status, "PENDING")
	err := transactions.ChangeStatus(testObject.ID, "ACTIVE")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "seller does not have enough stocks")
	checkTransaction, err := transactionRepo.GetByID(testObject.ID)
	assert.Nil(t, err)
	assert.Equal(t, checkTransaction.Status, "PENDING")

	err = stockRepo.AddStockToUser(testStock.ID, testSeller.ID, 12)
	assert.Nil(t, err)

	err = transactions.ChangeStatus(testObject.ID, "ACTIVE")
	assert.Nil(t, err)

	checkTransaction, err = transactionRepo.GetByID(testObject.ID)
	assert.Nil(t, err)
	assert.Equal(t, "ACTIVE", checkTransaction.Status)
}

func TestStatusSetPending(t *testing.T) {

	_, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testSeller := &models.User{Username: "Jack", Password: "abc123"}
	userRepo.Add(testSeller)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	stockRepo.Add(testStock)
	stockRepo.AddStockToUser(testStock.ID, testSeller.ID, 12)
	testObject := models.Transaction{
		SellerID: testSeller.ID,
		StockID:  testStock.ID,
		Quantity: 12,
		Type:     "SELL",
		Status:   "PENDING",
	}
	transactionRepo.Add(&testObject)
	assert.Equal(t, testObject.Status, "PENDING")
	err := transactions.ChangeStatus(testObject.ID, "ACTIVE")
	assert.Nil(t, err)
	err = transactions.ChangeStatus(testObject.ID, "PENDING")
	assert.Nil(t, err)

	checkTransaction, err := transactionRepo.GetByID(testObject.ID)
	assert.Nil(t, err)
	assert.Equal(t, "PENDING", checkTransaction.Status)
}

func TestStatusSetClosed(t *testing.T) {

	_, userRepo, stockRepo, transactionRepo := api.TestAssemblyRoute()
	testSeller := &models.User{Username: "Jack", Password: "abc123", Email: "aba", Funds: 0}
	err := userRepo.Add(testSeller)
	assert.Nil(t, err)
	testBuyer := &models.User{Username: "Jacka", Password: "abc123", Funds: 1000, Email: "abaa"}
	err = userRepo.Add(testBuyer)
	assert.Nil(t, err)
	testStock := &models.Stock{Name: "Stock1", Quantity: 123}
	err = stockRepo.Add(testStock)
	assert.Nil(t, err)
	err = stockRepo.AddStockToUser(testStock.ID, testSeller.ID, 123)
	assert.Nil(t, err)
	testObject := models.Transaction{
		SellerID: testSeller.ID,
		BuyerID:  testBuyer.ID,
		StockID:  testStock.ID,
		Quantity: 12,
		Type:     "SELL",
		Status:   "PENDING",
		Price:    1,
	}
	transactionRepo.Add(&testObject)
	assert.Equal(t, testObject.Status, "PENDING")
	err = transactions.ChangeStatus(testObject.ID, "ACTIVE")
	assert.Nil(t, err)
	err = transactions.ChangeStatus(testObject.ID, "CLOSED")
	assert.Nil(t, err)

	checkTransaction, err := transactionRepo.GetByID(testObject.ID)
	assert.Nil(t, err)
	assert.Equal(t, "CLOSED", checkTransaction.Status)
}

func TestTransactionCosts(t *testing.T) {
	_, _, _, transactionRepo := api.TestAssemblyRoute()
	testObject := models.Transaction{
		Quantity: 12,
		Type:     "SELL",
		Status:   "PENDING",
		Price:    2,
	}
	transactionRepo.Add(&testObject)
	cost, err := transactionRepo.TransactionCostsbyID(testObject.ID)
	assert.Nil(t, err)
	assert.Equal(t, 24.0, cost)
}
