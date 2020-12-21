package handler

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	//ErrRedirectNotFound "Redirect not found" error
	ErrRedirectNotFound = errors.New("Redirect not found")
	//ErrRedirectInvalid "Redirect invalid" error
	ErrRedirectInvalid = errors.New("Redirect invalid")
)

type redirectService struct {
	redirectRepo RedirectRepository
}

//NewRedirectService returns redirect service
func NewRedirectService(redirectRepo RedirectRepository) RedirectService {

	return &redirectService{
		redirectRepo,
	}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepo.Find(code)
}

func (r *redirectService) Store(redirect *Redirect) error {

	err := validate.Validate(redirect)

	if err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirection")
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()

	return r.redirectRepo.Store(redirect)
}
