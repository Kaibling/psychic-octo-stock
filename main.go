package main

import (
	"log"
	"net/http"

	"github.com/Kaibling/psychic-octo-stock/api"
)

func main() {
	r := api.AssembleServer()
	log.Println("=== Started ===")
	http.ListenAndServe("0.0.0.0:8080", r)
}
