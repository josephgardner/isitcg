package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/josephgardner/isitcg/internal/isitcg"
)

// getClientIP extracts the client IP from the request,
// checking X-Forwarded-For header first (for proxies like Dokku/nginx)
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by reverse proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip := r.RemoteAddr
	// Strip port if present
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

const (
	ROUTE_VIEW = "view"
)

func router(ingredientHandler isitcg.IngredientHandler, renders renders, counter isitcg.Counter) *mux.Router {

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
				return
			}

			productName := r.PostFormValue("productname")
			ingredients := r.PostFormValue("ingredients")

			hash := ingredientHandler.CreateHash(productName, ingredients)

			// Count on POST (not GET) to avoid bot traffic
			product := ingredientHandler.ProductFromHash(hash)
			res := ingredientHandler.ResultsFromProduct(product)
			clientIP := getClientIP(r)
			counter.Count(context.Background(), product, res, clientIP)

			if url, err := router.Get(ROUTE_VIEW).URL("hash", hash); err != nil {
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
			hash := mux.Vars(r)["hash"]
			product := ingredientHandler.ProductFromHash(hash)
			res := ingredientHandler.ResultsFromProduct(product)
			res.Hash = hash
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
