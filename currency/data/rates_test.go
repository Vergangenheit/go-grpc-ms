package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestExchangeRates(t *testing.T) {
	tr, err := NewRates(hclog.Default())
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Rates %#v", tr.rate)

}
