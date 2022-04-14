package handlers

import (
	"net/http"

	"github.com/Vergangenheit/go-grpc-ms/product-api/data"
)

func (*Products) Delete(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)

	p.l.Println("[DEBUG] deleting record id", id)

	err := data.DeleteProduct()

}
