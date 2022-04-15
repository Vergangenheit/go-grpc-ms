package handlers

import (
	"net/http"

	"github.com/Vergangenheit/go-grpc-ms/product-api/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products from the database
// responses:
//  200: productsResponse

// ListAll handles GET request and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET Products")
	prods := data.GetProducts()
	err := data.ToJSON(prods, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}
