package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/josephgardner/isitcg/internal/isitcg"
)

const (
	ROUTE_VIEW = "view"
)

func router(ingredientHandler isitcg.IngredientHandler, renders renders) *mux.Router {

	router := mux.NewRouter()

	router.NewRoute().
		Path("/").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			renders.Index(w, isitcg.Product{})
		})

	router.NewRoute().
		Path("/").
		Methods(http.MethodPost).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if url, err := router.Get(ROUTE_VIEW).URL(
				"hash",
				ingredientHandler.CreateHash(
					r.PostFormValue("productname"),
					r.PostFormValue("ingredients"),
				),
			); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				http.Redirect(w, r, url.String(), http.StatusSeeOther)
			}
		})

	router.NewRoute().
		Name(ROUTE_VIEW).
		Path("/view/{hash}").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := ingredientHandler.ResultsFromHash(mux.Vars(r)["hash"])
			renders.Results(w, res)
		})

	router.NewRoute().
		Path("/edit/{hash}").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			res := ingredientHandler.ProductFromHash(mux.Vars(r)["hash"])
			renders.Index(w, res)
		})
	return router
}
