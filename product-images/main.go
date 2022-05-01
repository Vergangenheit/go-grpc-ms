package main

import (
	"github.com/Vergangenheit/go-grpc-ms/product-images/files"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":8080", "Bind address for the server")
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
	stor, err := files.NewLocal()
}
