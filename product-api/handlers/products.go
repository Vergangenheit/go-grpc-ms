// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Vergangenheit/go-grpc-ms/product-api/data"
	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST product")
	// create muy new prodcut object
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	p.l.Println("Handle PUT product", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err2 := data.UpdateProduct(id, &prod)
	if err2 == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err2 != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

// define a key for the context
type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// validate the request deserialize product
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		// validate the product
		valid_err := prod.Validate()
		if valid_err != nil {
			p.l.Println("[ERROR] validating product", valid_err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", valid_err),
				http.StatusBadRequest,
			)
			return
		}

		// put the product in the request context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		// create a copy of that request
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
