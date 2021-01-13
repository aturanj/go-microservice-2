package shortener

//RedirectService is using for url-code redirection
type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
