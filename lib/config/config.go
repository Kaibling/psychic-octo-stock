package config

import (
	"log"
	"os"
)

var OSPREFIX = "POS"
var Config *Configuration

type Configuration struct {
	DBUrl       string
	Env         string
	LogFormat   string
	Currency    string
}

func NewConfig() *Configuration {
	// path, err := os.Getwd()
	// if err != nil {
	// 	log.Println(err)
	// }
	// dbURL := filepath.Join(path, "local.db")
	dbURL := getEnv("DBURL", "file::memory:")
	env := getEnv("ENV", "DEV")
	logFormat := getEnv("LOGFORMAT", "JSON")
	currency := getEnv("CURRENCY", "EUR")

	Config = &Configuration{DBUrl: dbURL, Env: env, LogFormat: logFormat, Currency: currency}
	return Config
}

func (s *Configuration) LogEnv() {
	log.Println("DBURL: " + s.DBUrl)
}

func getEnv(key string, defaultValue string) string {
	fullKey := OSPREFIX + "_" + key
	val := os.Getenv(OSPREFIX + "_" + key)
	if val == "" {
		if defaultValue != "" {
			return defaultValue
		}
		panic(fullKey + " is not set")
	}
	return val

}
