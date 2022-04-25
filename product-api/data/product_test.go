package data

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductMissingNameReturnsErr(t *testing.T) {
	// instantiate Product
	p := Product{
		ID:          45,
		Description: "a test product",
		Price:       23.5,
		SKU:         "ert-tam-mur",
	}
	// instantiate a Validation
	v := NewValidation()
	// validate prod
	val_errs := v.Validate(p).Errors()
	// assert val_errs is not empty(because it contains a validation error related to missing name in p)
	assert.Len(t, val_errs, 1)

}

func TestProductMissingPriceReturnsErr(t *testing.T) {
	// instantiate Product
	p := Product{
		ID:          45,
		Name:        "abc",
		Description: "a test product",
		SKU:         "ert-tam-mur",
	}
	// instantiate a Validation
	v := NewValidation()
	// validate prod
	val_errs := v.Validate(p).Errors()
	// assert val_errs are not empty
	assert.Len(t, val_errs, 1)
}

func TestProductInvalidSKUReturnsErr(t *testing.T) {
	p := Product{
		ID:          45,
		Name:        "Test Product",
		Description: "a test product",
		Price:       23.5,
		SKU:         "ert-222-mur",
	}
	// as above
	v := NewValidation()
	val_errs := v.Validate(p).Errors()
	assert.Len(t, val_errs, 1)
}

func TestValidProductDoesNOTReturnsErr(t *testing.T) {
	p := Product{
		ID:          45,
		Name:        "Test Product",
		Description: "a test product",
		Price:       23.5,
		SKU:         "ert-asd-mur",
	}
	v := NewValidation()
	err := v.Validate(p).Errors()
	assert.Len(t, err, 0)

}

func TestProductsToJSON(t *testing.T) {
	p := Product{
		ID:          45,
		Name:        "Test Product",
		Description: "a test product",
		Price:       23.5,
		SKU:         "ert-asd-mur",
	}
	b := bytes.NewBufferString("")
	err := ToJSON(p, b)
	assert.NoError(t, err)

}
