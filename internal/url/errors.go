package url

import "errors"

var (
	ErrAlreadyExists         = errors.New("URL already exists")
	ErrNotFound              = errors.New("URL not found")
	ErrAliasGenerationFailed = errors.New("alias generation failed")
)
