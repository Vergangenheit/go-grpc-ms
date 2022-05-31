package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	protos "github.com/Vergangenheit/go-grpc-ms/currency/currency"
	"github.com/Vergangenheit/go-grpc-ms/product-api/data"
	"github.com/Vergangenheit/go-grpc-ms/product-api/handlers"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func main() {
	// reference to the handler
	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// create client
	var cc protos.CurrencyClient = protos.NewCurrencyClient(conn)
	// refer to the handler
	ph := handlers.NewProducts(l, v, cc)
	// create servemux
	sm := mux.NewRouter()
	// create router for put request
	getRouter := sm.Methods("GET").Subrouter()
	getRouter.HandleFunc("/products", ph.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)
	// creater router for put requests
	putRouter := sm.Methods("PUT").Subrouter()
	putRouter.HandleFunc("/products", ph.Update)
	putRouter.Use(ph.MiddlewareValidateProduct)
	// router for post requests
	postRouter := sm.Methods("POST").Subrouter()
	postRouter.HandleFunc("/products", ph.Create)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// specify options
	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

	// instantiate an http server
	s := &http.Server{
		Addr:         ":8080",
		Handler:      ch(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	go func() {
		l.Println("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}() //non blocking function

	sigChan := make(chan os.Signal)
	// the service signal notify will broadcast to the chan whenever there's a interrupt command or kill command
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("received terminate, graceful shutdown", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
