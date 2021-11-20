package main

import "github.com/Kaibling/psychic-octo-stock/api"

func main() {
	r := api.AssembleServer()
	r.Run()
}
