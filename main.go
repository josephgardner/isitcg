package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/josephgardner/isitcg/internal/isitcg"
)

func main() {
	rules, err := isitcg.LoadRules("ingredientrules.json")
	if err != nil {
		panic("failed to load rules")
	}
	isitcg.NewIngredientHandler(rules)
	router := router(renderer())

	wwwroot := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", wwwroot))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "5000"
	}

	log.Printf("isitcg listening at http://0.0.0.0:%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port),
		handlers.CombinedLoggingHandler(os.Stdout, router)))
}