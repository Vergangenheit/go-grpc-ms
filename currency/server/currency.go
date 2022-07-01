package server

import (
	"context"

	protos "github.com/Vergangenheit/go-grpc-ms/currency/currency"
	"github.com/Vergangenheit/go-grpc-ms/currency/data"
	hclog "github.com/hashicorp/go-hclog"
)

type Currency struct {
	rates *data.ExchangeRates
	log   hclog.Logger
}

//constructor
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	return &Currency{r, l}
}

// implement CurrencyServer interface
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		c.log.Error("Can't get rates", "err", err)
		return nil, err
	}
	// return message
	return &protos.RateResponse{Rate: rate}, nil
}
