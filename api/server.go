package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/Kaibling/psychic-octo-stock/apimiddleware"
	"github.com/Kaibling/psychic-octo-stock/lib/config"
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/lib/transmission"
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
	r.Use(NewStructuredLogger(logger))
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

	//log.SetFormatter(&log.JSONFormatter{})
	// log.SetFormatter(
	// 	&easy.Formatter{
	// 		TimestampFormat: "2006-01-02 15:04:05",
	// 		LogFormat:       "[%lvl%]: %time% - %msg%",
	// 	})

	// log.SetOutput(os.Stdout)
	// log.SetLevel(log.DebugLevel)
	logger := logrus.New()
	logger.SetFormatter(
		&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%",
		})

	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	logger.Formatter = &logrus.JSONFormatter{
		// disable, as we set our own
		DisableTimestamp: true,
	}
	return logger
}

func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *logrus.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	response := utility.GetContext("responseObject", r).(*transmission.Response)
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method
	logFields["request_id"] = response.GetRequestId()

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
