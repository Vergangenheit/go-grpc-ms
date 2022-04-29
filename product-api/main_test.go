package main

import (
	"fmt"
	"testing"

	"github.com/Vergangenheit/go-grpc-ms/product-api/sdk/client"
	"github.com/Vergangenheit/go-grpc-ms/product-api/sdk/client/products"
)

func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:8080")
	c := client.NewHTTPClientWithConfig(nil, cfg)
	params := products.NewListProductsParams()
	prods, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(prods)
}
