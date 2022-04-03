package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "nics",
		Price: 1.0,
		SKU:   "ams-abc-dem",
	}

	err := p.Validate()
	if err != nil {
		t.Fatal(err)
	}
}
