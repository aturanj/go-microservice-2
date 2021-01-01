package messagepack

import (
	shortener "go-url-shortener/handler"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

//Redirect struct
type Redirect struct{}

//Decode binary input
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {

	redirect := &shortener.Redirect{}

	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

//Encode url input
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {

	rawMsg, err := msgpack.Marshal(input)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}
