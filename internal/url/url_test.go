package url_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/kodeyeen/shortify/internal/domain"
	"github.com/kodeyeen/shortify/internal/dto"
	mockgen "github.com/kodeyeen/shortify/internal/generation/mock"
	"github.com/kodeyeen/shortify/internal/persistence"
	mockpers "github.com/kodeyeen/shortify/internal/persistence/mock"
	"github.com/kodeyeen/shortify/internal/url"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
	type Given struct {
		req *dto.CreateURLRequest

		url *domain.URL

		urlID  int64
		urlErr error

		alias    string
		aliasErr error
	}

	type Expected struct {
		svcResp *dto.CreateURLResponse
		svcErr  error
	}

	testCases := map[string]struct {
		given    Given
		expected Expected
	}{
		"Success": {
			Given{
				req: &dto.CreateURLRequest{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
				},

				url: &domain.URL{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "randomstri",
				},

				urlID:  1,
				urlErr: nil,

				alias:    "randomstri",
				aliasErr: nil,
			},
			Expected{
				svcResp: &dto.CreateURLResponse{
					ID:       1,
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "randomstri",
				},
				svcErr: nil,
			},
		},
		"URL already exists": {
			Given{
				req: &dto.CreateURLRequest{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
				},

				url: &domain.URL{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "randomstri",
				},

				urlID:  1,
				urlErr: persistence.ErrURLAlreadyExists,

				alias:    "randomstri",
				aliasErr: nil,
			},
			Expected{
				svcResp: nil,
				svcErr:  url.ErrAlreadyExists,
			},
		},
		"Alias generation failed": {
			Given{
				req: &dto.CreateURLRequest{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
				},

				url: &domain.URL{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "",
				},

				urlID:  0,
				urlErr: nil,

				alias:    "",
				aliasErr: io.ErrShortBuffer,
			},
			Expected{
				svcResp: nil,
				svcErr:  url.ErrAliasGenerationFailed,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			ctx := context.Background()

			aliases := mockgen.NewAliasProvider(t)
			aliases.On("Generate", ctx, tc.given.req.Original).
				Return(tc.given.alias, tc.given.aliasErr).
				Once()

			var urls *mockpers.URLRepository

			if tc.given.aliasErr == nil {
				urls = mockpers.NewURLRepository(t)
				urls.On("Add", ctx, tc.given.url).
					Return(tc.given.urlID, tc.given.urlErr).
					Once()
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			svc := url.NewService(urls, aliases, log)

			// When
			resp, err := svc.Create(ctx, tc.given.req)

			// Then
			require.Equal(t, tc.expected.svcResp, resp)
			require.ErrorIs(t, tc.expected.svcErr, err)
		})
	}
}

func TestService_GetByAlias(t *testing.T) {
	type Given struct {
		req *dto.GetURLByAliasRequest

		url    *domain.URL
		urlErr error
	}

	type Expected struct {
		svcResp *dto.GetURLByAliasResponse
		svcErr  error
	}

	testCases := map[string]struct {
		given    Given
		expected Expected
	}{
		"Success": {
			Given{
				req: &dto.GetURLByAliasRequest{
					Alias: "fjda89fadb",
				},

				url: &domain.URL{
					ID:       1,
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "fjda89fadb",
				},
				urlErr: nil,
			},
			Expected{
				svcResp: &dto.GetURLByAliasResponse{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "fjda89fadb",
				},
				svcErr: nil,
			},
		},
		"Not found": {
			Given{
				req: &dto.GetURLByAliasRequest{
					Alias: "fjda89fadb",
				},

				url:    nil,
				urlErr: persistence.ErrURLNotFound,
			},
			Expected{
				svcResp: nil,
				svcErr:  url.ErrNotFound,
			},
		},
		"Other error": {
			Given{
				req: &dto.GetURLByAliasRequest{
					Alias: "fjda89fadb",
				},

				url:    nil,
				urlErr: errors.New("some retrieval error"),
			},
			Expected{
				svcResp: nil,
				svcErr:  fmt.Errorf("failed to get URL by alias: %w", errors.New("some retrieval error")),
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			ctx := context.Background()

			aliases := mockgen.NewAliasProvider(t)

			urls := mockpers.NewURLRepository(t)
			urls.On("FindByAlias", ctx, tc.given.req.Alias).
				Return(tc.given.url, tc.given.urlErr).
				Once()

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			svc := url.NewService(urls, aliases, log)

			// When
			resp, err := svc.GetByAlias(ctx, tc.given.req)

			// Then
			require.Equal(t, tc.expected.svcResp, resp)
			require.Equal(t, tc.expected.svcErr, err)
		})
	}
}
