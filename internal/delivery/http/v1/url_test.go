package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	httpdel "github.com/kodeyeen/shortify/internal/delivery/http/v1"
	"github.com/kodeyeen/shortify/internal/dto"
	"github.com/kodeyeen/shortify/internal/url"
	"github.com/kodeyeen/shortify/internal/urlmock"
	"github.com/kodeyeen/shortify/v1"

	"github.com/stretchr/testify/require"
)

func TestURLController_Create(t *testing.T) {
	type Given struct {
		reqBody []byte

		svcReq  *dto.CreateURLRequest
		svcResp *dto.CreateURLResponse
		svcErr  error
	}

	type Expected struct {
		statusCode  int
		successResp *shortify.CreateURLResponse
		errResp     *shortify.ErrorResponse
	}

	testCases := map[string]struct {
		given    Given
		expected Expected
	}{
		"Success": {
			Given{
				reqBody: []byte(`{"original": "https://example.com/longlonglonglonglonglonglonglong"}`),

				svcReq: &dto.CreateURLRequest{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
				},
				svcResp: &dto.CreateURLResponse{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "https://shortify.com/shortshort",
				},
				svcErr: nil,
			},
			Expected{
				statusCode: http.StatusCreated,
				successResp: &shortify.CreateURLResponse{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "https://shortify.com/shortshort",
				},
				errResp: nil,
			},
		},
		"Invalid request body": {
			Given{
				reqBody: []byte(`not json`),

				svcReq:  nil,
				svcResp: nil,
				svcErr:  nil,
			},
			Expected{
				statusCode:  http.StatusBadRequest,
				successResp: nil,
				errResp: &shortify.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "Invalid request body",
				},
			},
		},
		"Empty Original": {
			Given{
				reqBody: []byte(`{"original": ""}`),

				svcReq:  nil,
				svcResp: nil,
				svcErr:  nil,
			},
			Expected{
				statusCode:  http.StatusBadRequest,
				successResp: nil,
				errResp: &shortify.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "Field 'original' is missing",
				},
			},
		},
		"Invalid Original": {
			Given{
				reqBody: []byte(`{"original": "invalidurl"}`),

				svcReq:  nil,
				svcResp: nil,
				svcErr:  nil,
			},
			Expected{
				statusCode:  http.StatusBadRequest,
				successResp: nil,
				errResp: &shortify.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "Field 'original' is not a valid URL",
				},
			},
		},
		"URL already exists": {
			Given{
				reqBody: []byte(`{"original": "https://example.com/longlonglonglonglonglonglonglong"}`),

				svcReq: &dto.CreateURLRequest{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
				},
				svcResp: nil,
				svcErr:  url.ErrAlreadyExists,
			},
			Expected{
				statusCode:  http.StatusConflict,
				successResp: nil,
				errResp: &shortify.ErrorResponse{
					Status:  http.StatusConflict,
					Message: "URL already exists",
				},
			},
		},
		"Other": {
			Given{
				reqBody: []byte(`{"original": "https://example.com/longlonglonglonglonglonglonglong"}`),

				svcReq: &dto.CreateURLRequest{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
				},
				svcResp: nil,
				svcErr:  errors.New("svc error"),
			},
			Expected{
				statusCode:  http.StatusInternalServerError,
				successResp: nil,
				errResp: &shortify.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: http.StatusText(http.StatusInternalServerError),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodPost, "/api/v1/urls", bytes.NewReader(tc.given.reqBody))
			require.NoError(t, err)

			ctx := req.Context()

			svc := urlmock.NewService(t)

			if tc.given.svcResp != nil || tc.given.svcErr != nil {
				svc.On("Create", ctx, tc.given.svcReq).
					Return(tc.given.svcResp, tc.given.svcErr).
					Once()
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			clr := httpdel.NewURLController(svc, log)

			// When
			clr.Create(rr, req)

			// Then
			require.Equal(t, tc.expected.statusCode, rr.Code)

			if tc.expected.errResp != nil {
				var resp shortify.ErrorResponse

				err = json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, tc.expected.errResp, &resp)
			} else {
				var resp shortify.CreateURLResponse

				err = json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, tc.expected.successResp, &resp)
			}
		})
	}
}

func TestURLController_GetByAlias(t *testing.T) {
	type Given struct {
		alias string

		svcReq  *dto.GetURLByAliasRequest
		svcResp *dto.GetURLByAliasResponse
		svcErr  error
	}

	type Expected struct {
		statusCode  int
		successResp *shortify.GetURLByAliasResponse
		errResp     *shortify.ErrorResponse
	}

	testCases := map[string]struct {
		given    Given
		expected Expected
	}{
		"Success": {
			Given{
				alias: "fjsido39jf",

				svcReq: &dto.GetURLByAliasRequest{
					Alias: "fjsido39jf",
				},
				svcResp: &dto.GetURLByAliasResponse{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "fjsido39jf",
				},
				svcErr: nil,
			},
			Expected{
				statusCode: http.StatusOK,
				successResp: &shortify.GetURLByAliasResponse{
					Original: "https://example.com/longlonglonglonglonglonglonglong",
					Alias:    "fjsido39jf",
				},
				errResp: nil,
			},
		},
		"Empty alias": {
			Given{
				alias: "",

				svcReq:  nil,
				svcResp: nil,
				svcErr:  nil,
			},
			Expected{
				statusCode:  http.StatusBadRequest,
				successResp: nil,
				errResp: &shortify.ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "Alias is empty",
				},
			},
		},
		"Not found": {
			Given{
				alias: "fjsido39jf",

				svcReq: &dto.GetURLByAliasRequest{
					Alias: "fjsido39jf",
				},
				svcResp: nil,
				svcErr:  url.ErrNotFound,
			},
			Expected{
				statusCode:  http.StatusNotFound,
				successResp: nil,
				errResp: &shortify.ErrorResponse{
					Status:  http.StatusNotFound,
					Message: http.StatusText(http.StatusNotFound),
				},
			},
		},
		"Other": {
			Given{
				alias: "fjsido39jf",

				svcReq: &dto.GetURLByAliasRequest{
					Alias: "fjsido39jf",
				},
				svcResp: nil,
				svcErr:  errors.New("svc error"),
			},
			Expected{
				statusCode:  http.StatusInternalServerError,
				successResp: nil,
				errResp: &shortify.ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: http.StatusText(http.StatusInternalServerError),
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Given
			rr := httptest.NewRecorder()

			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/urls/%s", tc.given.alias), nil)
			require.NoError(t, err)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("alias", tc.given.alias)

			ctx := context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
			req = req.WithContext(ctx)

			svc := urlmock.NewService(t)

			if tc.given.svcResp != nil || tc.given.svcErr != nil {
				svc.On("GetByAlias", ctx, tc.given.svcReq).
					Return(tc.given.svcResp, tc.given.svcErr).
					Once()
			}

			log := slog.New(slog.NewTextHandler(io.Discard, nil))

			clr := httpdel.NewURLController(svc, log)

			// When
			clr.GetByAlias(rr, req)

			// Then
			require.Equal(t, tc.expected.statusCode, rr.Code)

			if tc.expected.errResp != nil {
				var resp shortify.ErrorResponse

				err = json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, tc.expected.errResp, &resp)
			} else {
				var resp shortify.GetURLByAliasResponse

				err = json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, tc.expected.successResp, &resp)
			}
		})
	}
}
