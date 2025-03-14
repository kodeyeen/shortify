package persistence

import (
	"errors"
	"net/url"
)

var (
	ErrURLNotFound      = errors.New("URL not found")
	ErrURLAlreadyExists = errors.New("URL already exists")

	ErrDuplicateAlias = errors.New("duplicate alias")
)

// type URLRepository interface {
// 	Add(ctx context.Context, u *domain.URL) (int64, error)
// 	FindByAlias(ctx context.Context, alias string) (*domain.URL, error)
// }

func NewConnString(driver, username, password, host, db string) string {
	u := url.URL{
		Scheme: driver,
		User:   url.UserPassword(username, password),
		Host:   host,
		Path:   db,
	}

	return u.String()
}
