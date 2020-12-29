package json

import (
	"encoding/json"

	shortener "go-url-shortener/handler"

	"github.com/pkg/errors"
)

//Redirect struct
type Redirect struct{}

//Decode binary input
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {

	redirect := &shortener.Redirect{}

	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

//Encode url input
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {

	rawMsg, err := json.Marshal(input)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}
