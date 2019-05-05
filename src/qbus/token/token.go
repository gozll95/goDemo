package token

type Token interface {
	GetToken() (string, error)
}
