package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// interface
type Hello struct {
	l *log.Logger
}

//
func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

// method that satify the interface
func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello World")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "OOps", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "hello %s", d)
}
