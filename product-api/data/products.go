package data

import (
	"fmt"
	"time"
)

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for this product
	//
	// required: true
	// min: 1
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

// make the returned pointers to struct a type in itself to link methods to it
type Products []*Product

// add a validate method on ur type
// func (p *Product) Validate() error {
// 	validate := validator.New()
// 	validate.RegisterValidation("sku", validateSKU)
// 	return validate.Struct(p)
// }

// custom validation function for sku field
// func validateSKU(fl validator.FieldLevel) bool {
// 	// use a regex
// 	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
// 	matches := re.FindAllString(fl.Field().String(), -1)
// 	if len(matches) != 1 {
// 		return false
// 	}
// 	return true
// }

// data access model
func GetProducts() Products {
	return productList
}

func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

// store into database
func AddProduct(p Product) {
	// set product id
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, &p)
}

func UpdateProduct(p Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}
	// update product list
	productList[i] = &p

	return nil
}

func DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == 1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil

}

var ErrProductNotFound = fmt.Errorf("product not found")

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}
	return -1
}

func getNextID() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

// productList is a hard coded list of products for this
// example data source
var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
