package api

import (
	"bytes"
	"context"
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
	"github.com/lucsky/cuid"
)

func AssembleServer() *chi.Mux {
	configData := config.NewConfig()
	configData.LogEnv()
	db := database.NewDatabaseConnector(configData.DBUrl)
	db.Connect()
	migrateDB(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(injectData("hmacSecret", []byte(config.Config.TokenSecret)))
	initRepos(r, db)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	BuildRouter(r)
	displayRoutes(r)
	return r
}
func TestAssemblyRoute() (*chi.Mux, *repositories.UserRepository, *repositories.StockRepository, *repositories.TransactionRepository, func(r http.Handler, method, path string, jsonStr []byte) *httptest.ResponseRecorder) {
	configData := config.NewConfig()
	configData.LogEnv()
	db := database.NewDatabaseConnector(configData.DBUrl)
	db.Connect()
	token := migrateDB(db)

	userRepo := repositories.NewUserRepository(db)
	repositories.SetUserRepo(userRepo)
	stockRepo := repositories.NewStockRepository(db)
	repositories.SetStockRepo(stockRepo)
	transactionRepo := repositories.NewTransactionRepository(db)
	repositories.SetTransactionRepo(transactionRepo)

	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(injectData("userRepo", userRepo))
	r.Use(injectData("stockRepo", stockRepo))
	r.Use(injectData("transactionRepo", transactionRepo))
	r.Use(injectData("hmacSecret", []byte(config.Config.TokenSecret)))
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

func migrateDB(db database.DBConnector) string {

	db.Migrate(&models.User{})
	db.Migrate(&models.Stock{})
	db.Migrate(&models.StockToUser{})
	db.Migrate(&models.Transaction{})

	if err := db.FindByWhere(&models.User{}, "username = ?", []interface{}{"admin"}); err != nil {
		log.Println("no admin found. Creating one")
		password := cuid.New()
		admin := &models.User{Username: "admin", Password: password, Email: "admin@local"}
		db.Add(admin)
		token, _ := utility.GenerateToken(admin.Username, []byte(config.Config.TokenSecret))
		log.Printf("user: %s, password: %s\ntoken: %s\n", admin.Username, password, token)
		return token
	}
	return ""
}

func initRepos(r *chi.Mux, db database.DBConnector) {
	userRepo := repositories.NewUserRepository(db)
	repositories.SetUserRepo(userRepo)
	stockRepo := repositories.NewStockRepository(db)
	repositories.SetStockRepo(stockRepo)
	transactionRepo := repositories.NewTransactionRepository(db)
	repositories.SetTransactionRepo(transactionRepo)

	r.Use(injectData("userRepo", userRepo))
	r.Use(injectData("stockRepo", stockRepo))
	r.Use(injectData("transactionRepo", transactionRepo))
}

func displayRoutes(r *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.Printf("Logging err: %s\n", err.Error())
	}
}
