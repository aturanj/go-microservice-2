package handler

type RedirectSerializer interface {
	Decode(input []byte) (*Redirect, error)
}
