package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func main() {
	router := router(renderer())

	wwwroot := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", wwwroot))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "5000"
	}

	log.Printf("isitcg listening at %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port),
		handlers.CombinedLoggingHandler(os.Stdout, router)))
}
