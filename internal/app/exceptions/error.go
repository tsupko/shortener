package exceptions

import "errors"

var (
	ErrURLAlreadyExist  = errors.New("url already exists")
	ErrHashAlreadyExist = errors.New("hash already exists")
	ErrURLNotFound      = errors.New("url not found")
)
