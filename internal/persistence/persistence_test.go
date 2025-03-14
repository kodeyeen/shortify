package persistence_test

import (
	"testing"

	"github.com/kodeyeen/shortify/internal/persistence"
	"github.com/stretchr/testify/require"
)

func TestNewConnString(t *testing.T) {
	type Given struct {
		driver   string
		username string
		password string
		host     string
		db       string
	}

	type Expected struct {
		res string
	}

	testCases := map[string]struct {
		given    Given
		expected Expected
	}{
		"Success": {
			Given{
				driver:   "postgres",
				username: "username",
				password: "password",
				host:     "host",
				db:       "db",
			},
			Expected{
				res: "postgres://username:password@host/db",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			connString := persistence.NewConnString(
				tc.given.driver,
				tc.given.username,
				tc.given.password,
				tc.given.host,
				tc.given.db,
			)

			require.Equal(t, tc.expected.res, connString)
		})
	}
}
