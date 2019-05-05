package httpclient

import "errors"

var (
	ErrResponseRead  = errors.New("Response body has been read.")
	ErrInvalidParams = errors.New("Request params is invalid")
)
