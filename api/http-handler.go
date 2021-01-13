package api

import (
	"go-url-shortener/shortener"
	"net/http"
)

//RedirectHandler interface
type RedirectHandler interface {
	Get(http.ResponseWriter, http.Request)
	Post(http.ResponseWriter, http.Request)
}

type handler struct {
	redirectService shortener.RedirectService
}

//NewHandler returns redirect service
func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService: redirectService}
}
