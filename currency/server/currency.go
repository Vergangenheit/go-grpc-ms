package server

import (
	"context"

	protos "github.com/Vergangenheit/go-grpc-ms/currency/currency"
	hclog "github.com/hashicorp/go-hclog"
)

type Currency struct {
	log hclog.Logger
}

//constructor
func NewCurrency(l hclog.Logger) *Currency {
	return &Currency{log: l}
}

// implement CurrencyServer interface
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	// return message
	return &protos.RateResponse{Rate: 0.5}, nil
}
