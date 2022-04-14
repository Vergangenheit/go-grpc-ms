package handlers

import (
	"net/http"

	"github.com/Vergangenheit/go-grpc-ms/product-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products from the database
// responses:
//  200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
