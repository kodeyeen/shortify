package http

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/kodeyeen/shortify/internal/dto"
	"github.com/kodeyeen/shortify/internal/url"
	"github.com/kodeyeen/shortify/v1"
)

type URLService interface {
	Create(ctx context.Context, req *dto.CreateURLRequest) (*dto.CreateURLResponse, error)
	GetByAlias(ctx context.Context, req *dto.GetURLByAliasRequest) (*dto.GetURLByAliasResponse, error)
}

type URLController struct {
	urls URLService

	log *slog.Logger
}

func NewURLController(urls URLService, log *slog.Logger) *URLController {
	return &URLController{
		urls: urls,

		log: log,
	}
}

// Create creates new URL and generates an alias for it
//
//	@Summary		Create a URL
//	@Description	Create creates new URL and generates an alias for it
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			URL	body		shortify.CreateURLRequest	true	"Create URL"
//	@Success		200	{object}	shortify.CreateURLResponse
//	@Failure		400	{object}	shortify.ErrorResponse
//	@Failure		404	{object}	shortify.ErrorResponse
//	@Failure		500	{object}	shortify.ErrorResponse
//	@Router			/api/v1/urls [post]
func (c *URLController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := c.log.With(
		slog.String("handler", "Create"),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	var req shortify.CreateURLRequest

	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		log.Error("failed to decode request body", slog.String("error", err.Error()))

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, shortify.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid request body",
		})
		return
	}

	log.Info("request body decoded", slog.Any("request", req))

	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", slog.String("error", err.Error()))

		validatorErrs := err.(validator.ValidationErrors)

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, shortify.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: formatErrs(validatorErrs),
		})
		return
	}

	out, err := c.urls.Create(ctx, &dto.CreateURLRequest{Original: req.Original})
	if err != nil {
		if errors.Is(err, url.ErrAlreadyExists) {
			log.Info("url already exists", slog.String("url", req.Original))

			render.Status(r, http.StatusConflict)
			render.JSON(w, r, shortify.ErrorResponse{
				Status:  http.StatusConflict,
				Message: "URL already exists",
			})
			return
		}

		log.Info("failed to create URL", slog.String("error", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, shortify.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	log.Info("URL created", slog.Int64("id", out.ID))

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, shortify.CreateURLResponse{
		Original: out.Original,
		Alias:    out.Alias,
	})
}

// GetByAlias gets URL by its alias
//
//	@Summary		Get URL by its alias
//	@Description	Get URL by its alias
//	@Tags			urls
//	@Accept			json
//	@Produce		json
//	@Param			alias	path		string	true	"Get URL by alias"
//	@Success		200		{object}	shortify.GetURLByAliasResponse
//	@Failure		400		{object}	shortify.ErrorResponse
//	@Failure		404		{object}	shortify.ErrorResponse
//	@Failure		500		{object}	shortify.ErrorResponse
//	@Router			/api/v1/urls/{alias} [get]
func (c *URLController) GetByAlias(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := c.log.With(
		slog.String("handler", "GetByAlias"),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)

	alias := chi.URLParam(r, "alias")
	if alias == "" {
		log.Info("alias is empty")

		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, shortify.ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Alias is empty",
		})
		return
	}

	out, err := c.urls.GetByAlias(ctx, &dto.GetURLByAliasRequest{
		Alias: alias,
	})
	if err != nil {
		if errors.Is(err, url.ErrNotFound) {
			log.Info("URL not found", "alias", alias)

			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, shortify.ErrorResponse{
				Status:  http.StatusNotFound,
				Message: http.StatusText(http.StatusNotFound),
			})
			return
		}

		log.Error("failed to get URL by alias", slog.String("error", err.Error()))

		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, shortify.ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	log.Info("got URL by alias", slog.String("url", out.Original))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, shortify.GetURLByAliasResponse{
		Original: out.Original,
		Alias:    out.Alias,
	})
	// http.Redirect(w, r, resp.Original, http.StatusFound)
}
