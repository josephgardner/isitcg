package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	ROUTE_VIEW = "view"
)

func router(renders renders) *mux.Router {

	router := mux.NewRouter()

	router.NewRoute().
		Path("/").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			renders.Index(w, product{})
		})

	router.NewRoute().
		Path("/").
		Methods(http.MethodPost).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url, err := router.Get(ROUTE_VIEW).URL("hash", "xyz123")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			http.Redirect(w, r, url.String(), http.StatusSeeOther)
		})

	router.NewRoute().
		Name(ROUTE_VIEW).
		Path("/view/{hash}").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			renders.Results(w, results{})
		})

	router.NewRoute().
		Path("/edit/{hash}").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			renders.Index(w, product{"hello", "edit"})
		})
	return router
}
