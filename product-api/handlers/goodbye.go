package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Goodbye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	g.l.Println("Goodye World")

	d, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, "OOps", http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, "Goodbye %s", d)
}
