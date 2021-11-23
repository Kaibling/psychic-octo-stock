package api

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/Kaibling/psychic-octo-stock/lib/config"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/lib/utility"
	"github.com/Kaibling/psychic-octo-stock/models"
	"github.com/Kaibling/psychic-octo-stock/repositories"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func AssembleServer() *chi.Mux {
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
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(injectData("userRepo", userRepo))
	r.Use(injectData("stockRepo", stockRepo))
	r.Use(injectData("transactionRepo", transactionRepo))
	r.Use(injectData("hmacSecret", []byte("asdassasdsdsdswew")))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	BuildRouter(r)
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
	return r
}
func TestAssemblyRoute() (*chi.Mux, *repositories.UserRepository, *repositories.StockRepository, *repositories.TransactionRepository, func(r http.Handler, method, path string, jsonStr []byte) *httptest.ResponseRecorder) {
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
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(injectData("userRepo", userRepo))
	r.Use(injectData("stockRepo", stockRepo))
	r.Use(injectData("transactionRepo", transactionRepo))
	r.Use(injectData("hmacSecret", []byte("asdassasdsdsdswew")))
	token, _ := utility.GenerateToken(testUser.Username, []byte("asdassasdsdsdswew"))
	BuildRouter(r)
	PerformTestRequest := func(r http.Handler, method, path string, jsonStr []byte) *httptest.ResponseRecorder {

		req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}

	return r, userRepo, stockRepo, transactionRepo, PerformTestRequest
}

func injectData(key string, data interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), key, data)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
