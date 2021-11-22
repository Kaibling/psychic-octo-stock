package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/Kaibling/psychic-octo-stock/lib/config"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/gin-gonic/gin"
)

func AssembleServer() *gin.Engine {
	configData := config.NewConfig()
	configData.LogEnv()
	sdb := database.NewDatabaseConnector(configData.DBUrl)
	sdb.Connect()
	sdb.Migrate(&models.User{})
	sdb.Migrate(&models.Stock{})
	sdb.Migrate(&models.StockToUser{})
	sdb.Migrate(&models.Transaction{})

	userRepo := repositories.NewUserRepository(sdb)
	repositories.SetUserRepo(userRepo)
	stockRepo := repositories.NewStockRepository(sdb)
	repositories.SetStockRepo(stockRepo)
	transactionRepo := repositories.NewTransactionRepository(sdb)
	repositories.SetTransactionRepo(transactionRepo)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(injectData("userRepo", userRepo))
	r.Use(injectData("stockRepo", stockRepo))
	r.Use(injectData("transactionRepo", transactionRepo))
	r.Use(injectData("hmacSecret", []byte("asdassasdsdsdswew")))
	BuildRouter(r)
	return r
}
func TestAssemblyRoute() (*gin.Engine, *repositories.UserRepository, *repositories.StockRepository, *repositories.TransactionRepository, func(r http.Handler, method, path string, jsonStr []byte) *httptest.ResponseRecorder) {
	configData := config.NewConfig()
	configData.LogEnv()
	sdb := database.NewDatabaseConnector(configData.DBUrl)
	sdb.Connect()
	sdb.Migrate(&models.User{})
	sdb.Migrate(&models.Stock{})
	sdb.Migrate(&models.StockToUser{})
	sdb.Migrate(&models.Transaction{})

	userRepo := repositories.NewUserRepository(sdb)
	testUser := &models.User{Username: "testUser", Password: "testpassword"}
	userRepo.Add(testUser)

	repositories.SetUserRepo(userRepo)
	stockRepo := repositories.NewStockRepository(sdb)
	repositories.SetStockRepo(stockRepo)
	transactionRepo := repositories.NewTransactionRepository(sdb)
	repositories.SetTransactionRepo(transactionRepo)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(injectData("userRepo", userRepo))
	r.Use(injectData("stockRepo", stockRepo))
	r.Use(injectData("transactionRepo", transactionRepo))
	r.Use(injectData("hmacSecret", []byte("asdassasdsdsdswew")))
	token, _ := utility.GenerateToken(testUser.Username, []byte("asdassasdsdsdswew"))
	BuildRouter(r)
	PerformTestRequest := func(r http.Handler, method, path string, jsonStr []byte) *httptest.ResponseRecorder {

		req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	return r, userRepo, stockRepo, transactionRepo, PerformTestRequest
}

func injectData(key string, data interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(key, data)
		c.Next()
	}
}
