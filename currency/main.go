package main

import (
	"net"
	"os"

	protos "github.com/Vergangenheit/go-grpc-ms/currency/currency"
	"github.com/Vergangenheit/go-grpc-ms/currency/server"
	hclog "github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log := hclog.Default()
	// create grpc server
	gs := grpc.NewServer()
	cs := server.NewCurrency(log)
	// register grpc server
	protos.RegisterCurrencyServer(gs, cs)
	// enable reflection api
	reflection.Register(gs)
	l, err := net.Listen("tcp", ":9092")
	if err != nil {
		log.Error("Unable to listen", "error", err)
		os.Exit(1)
	}
	gs.Serve(l)
}
