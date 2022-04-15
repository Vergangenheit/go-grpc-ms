package data

import (
	"encoding/json"
	"io"
)

func ToJSON(i interface{}, w io.Writer) error {
	// encoder
	e := json.NewEncoder(w)
	return e.Encode(i)
}

func FromJSON(i interface{}, r io.Reader) error {
	// create and encoder
	e := json.NewDecoder(r)
	return e.Decode(i)
}
