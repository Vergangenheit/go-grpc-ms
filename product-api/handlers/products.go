package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/Vergangenheit/go-grpc-ms/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// we want to check what kind of request it is
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)
		// expect the id in the uri
		// decode this
		reg := regexp.MustCompile(`/([0-9]+)`)
		groups := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(groups) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(groups[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		idString := groups[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Unable to convert to number", id)
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.updateProducts(id, rw, r)
		return
	}
	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)

	// handle an update

}

//
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("handle POST product")
	// create muy new prodcut object
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	// store that piece of info
	p.l.Printf("Prod: %v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT product")

	prod := &data.Product{}

	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	err2 := data.UpdateProduct(id, prod)
	if err2 == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}
	if err2 != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
