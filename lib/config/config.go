package config

import (
	"log"
	"os"
)

var OSPREFIX = "POS"
var Config *Configuration

type Configuration struct {
	DBUrl string
}

func NewConfig() *Configuration {
	dbURL := os.Getenv(OSPREFIX + "_DBURL")
	if dbURL == "" {
		// path, err := os.Getwd()
		// if err != nil {
		// 	log.Println(err)
		// }
		dbURL = "file::memory:" // filepath.Join(path, "local.db")
	}
	Config = &Configuration{DBUrl: dbURL}
	return &Configuration{DBUrl: dbURL}
}

func (s *Configuration) LogEnv() {
	log.Println("DBURL: " + s.DBUrl)
}
