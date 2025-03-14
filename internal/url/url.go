package url

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/kodeyeen/shortify/internal/domain"
	"github.com/kodeyeen/shortify/internal/dto"
	"github.com/kodeyeen/shortify/internal/persistence"
)

type Repository interface {
	Add(ctx context.Context, u *domain.URL) (int64, error)
	FindByAlias(ctx context.Context, alias string) (*domain.URL, error)
}

type AliasProvider interface {
	Generate(ctx context.Context, original string) (string, error)
}

type Service struct {
	urls    Repository
	aliases AliasProvider

	log *slog.Logger
}

func NewService(urls Repository, aliases AliasProvider, log *slog.Logger) *Service {
	return &Service{
		urls:    urls,
		aliases: aliases,

		log: log,
	}
}

// Create creates new URL
func (s *Service) Create(ctx context.Context, req *dto.CreateURLRequest) (*dto.CreateURLResponse, error) {
	alias, err := s.aliases.Generate(ctx, req.Original)
	if err != nil {
		return nil, ErrAliasGenerationFailed
	}

	u := &domain.URL{
		Original: req.Original,
		Alias:    alias,
	}

	var id int64

	for {
		id, err = s.urls.Add(ctx, u)
		if err != nil {
			if errors.Is(err, persistence.ErrDuplicateAlias) {
				continue
			} else if errors.Is(err, persistence.ErrURLAlreadyExists) {
				return nil, ErrAlreadyExists
			}

			return nil, fmt.Errorf("failed to create URL: %w", err)
		}

		break
	}

	u.ID = id

	return &dto.CreateURLResponse{
		ID:       u.ID,
		Original: u.Original,
		Alias:    u.Alias,
	}, nil
}

// GetByAlias gets URL by its alias
func (s *Service) GetByAlias(ctx context.Context, req *dto.GetURLByAliasRequest) (*dto.GetURLByAliasResponse, error) {
	u, err := s.urls.FindByAlias(ctx, req.Alias)
	if err != nil {
		if errors.Is(err, persistence.ErrURLNotFound) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("failed to get URL by alias: %w", err)
	}

	return &dto.GetURLByAliasResponse{
		Original: u.Original,
		Alias:    u.Alias,
	}, nil
}
