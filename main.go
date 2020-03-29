package main

import (
	"log"
	"net/http"

	"github.com/ashrap/GO_test/server"
)

func main() {

	r := server.InitServer()

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8888", r))
}
