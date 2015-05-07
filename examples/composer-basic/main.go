package main

import (
	"github.com/bernos/composer"
	"log"
	"net/http"
)

func main() {
	server := composer.NewComposerHttpServer()
	log.Fatal(http.ListenAndServe(":1080", server))
}
