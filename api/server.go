package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/Kaibling/psychic-octo-stock/apimiddleware"
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
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	//easy "github.com/t-tomalak/logrus-easy-formatter"
)

func baseServer() (*chi.Mux, database.DBConnector, string) {
	configData := config.NewConfig()
	configData.LogEnv()
	logger := initLogging()

	db := database.NewDatabaseConnector(configData.DBUrl)
	db.Connect()
	token := migrateDB(db)

	r := chi.NewRouter()

	r.Use(apimiddleware.Response)
	r.Use(utility.NewStructuredLogger(logger))
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(injectData("hmacSecret", []byte(config.Config.TokenSecret)))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	return r, db, token
}

func AssembleServer() *chi.Mux {
	r, db, _ := baseServer()
	initRepos(r, db)

	BuildRouter(r)
	displayRoutes(r)
	return r
}
func TestAssemblyRoute() (*chi.Mux, map[string]interface{}, func(r http.Handler, method, path string, jsonStr []byte, headers *map[string]string) *httptest.ResponseRecorder) {
	r, db, token := baseServer()
	repos := initRepos(r, db)
	PerformTestRequest := func(r http.Handler, method, path string, jsonStr []byte, headers *map[string]string) *httptest.ResponseRecorder {

		req, _ := http.NewRequest(method, path, bytes.NewBuffer(jsonStr))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		if headers != nil {
			for k, v := range *headers {
				req.Header.Set(k, v)
			}
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	BuildRouter(r)
	return r, repos, PerformTestRequest
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
		fmt.Println("no admin found. Creating one")
		password := cuid.New()
		admin := &models.User{Username: "admin", Password: password, Email: "admin@local"}
		db.Add(admin)
		token, _ := utility.GenerateToken(admin.Username, []byte(config.Config.TokenSecret))
		fmt.Printf("user: %s, password: %s\ntoken: %s\n", admin.Username, password, token)
		return token
	}
	return ""
}

func initRepos(r *chi.Mux, db database.DBConnector) map[string]interface{} {
	repos := map[string]interface{}{}
	userRepo := repositories.NewUserRepository(db)
	repositories.SetUserRepo(userRepo)
	repos["userRepo"] = userRepo
	stockRepo := repositories.NewStockRepository(db)
	repositories.SetStockRepo(stockRepo)
	repos["stockRepo"] = stockRepo
	transactionRepo := repositories.NewTransactionRepository(db)
	repositories.SetTransactionRepo(transactionRepo)
	repos["transactionRepo"] = transactionRepo

	r.Use(injectData("userRepo", userRepo))
	r.Use(injectData("stockRepo", stockRepo))
	r.Use(injectData("transactionRepo", transactionRepo))
	return repos

}

func displayRoutes(r *chi.Mux) {

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
}

func initLogging() *logrus.Logger {

	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	if config.Config.LogFormat == "STRING" {
		logger.SetFormatter(
			&easy.Formatter{
				TimestampFormat: "2006-01-02 15:04:05",
				LogFormat:       "[%lvl%]: %time% - %req_id% %remote_addr% %user_agent%  %msg%\n",
			})

	} else if config.Config.LogFormat == "JSON" {
		logger.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: true,
		}
	}
	return logger
}
