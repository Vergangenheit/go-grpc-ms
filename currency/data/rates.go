package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	hclog "github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log  hclog.Logger
	rate map[string]float64
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rate: map[string]float64{}}
	err := er.getRates()
	if err != nil {
		er.log.Error("We got an error", "err", err)
		return er, err
	}
	return er, nil
}

//
func (e *ExchangeRates) GetRate(base string, dest string) (float64, error) {
	br, ok := e.rate[base]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", base)
	}
	dr, ok := e.rate[dest]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", dest)
	}

	return dr / br, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		e.log.Error("Can't get data from url ", "err", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected error code 200, got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// process xml data
	md := &Cubes{}
	err_d := xml.NewDecoder(resp.Body).Decode(&md)
	if err_d != nil {
		e.log.Error("Can't decode data ", "err", err_d)
	}
	// loop over collection, convert to float and put into map
	for _, m := range md.CubeData {
		rate, err := strconv.ParseFloat(m.Rate, 64)
		if err != nil {
			return err
		}
		e.rate[m.Currency] = rate
	}
	e.rate["EUR"] = 1
	return nil
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
