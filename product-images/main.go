package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Vergangenheit/go-grpc-ms/product-images/files"
	"github.com/Vergangenheit/go-grpc-ms/product-images/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":8081", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-images",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	// create a logger for the server from the default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// max filesize 5MB
	stor, err := files.NewLocal(*basePath, 1024*1000*5)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// create the handlers
	fh := handlers.NewFiles(stor, l)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// filename regex: {filename:[a-zA-Z]+\\.[a-z]{3}}
	ph := sm.Methods(http.MethodPost).Subrouter()
	ph.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.UploadRest)
	ph.HandleFunc("/", fh.UploadMultiPart)

	// get files
	gh := sm.Methods(http.MethodGet).Subrouter()
	gh.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))))

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))
	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to writer response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// Start the server
	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start the server", "error", err)
			os.Exit(1)
		}

	}()

	sigChan := make(chan os.Signal)
	// the service signal notify will broadcast to the chan whenever there's a interrupt command or kill command
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	l.Info("received terminate, graceful shutdown", "signal", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
