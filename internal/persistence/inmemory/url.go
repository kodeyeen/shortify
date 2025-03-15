package inmemory

import (
	"context"
	"sync"

	"github.com/kodeyeen/shortify/internal/domain"
	"github.com/kodeyeen/shortify/internal/persistence"
)

type URLRepository struct {
	originalIdx map[string]*domain.URL
	aliasIdx    map[string]*domain.URL

	mu *sync.RWMutex
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		originalIdx: map[string]*domain.URL{},
		aliasIdx:    map[string]*domain.URL{},

		mu: &sync.RWMutex{},
	}
}

func (r *URLRepository) Add(ctx context.Context, u *domain.URL) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.originalIdx[u.Original]; ok {
		return 0, persistence.ErrURLAlreadyExists
	}

	if _, ok := r.aliasIdx[u.Alias]; ok {
		return 0, persistence.ErrDuplicateAlias
	}

	id := int64(len(r.originalIdx) + 1)

	r.originalIdx[u.Original] = u
	r.aliasIdx[u.Alias] = u

	return id, nil
}

func (r *URLRepository) FindByAlias(ctx context.Context, alias string) (*domain.URL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.aliasIdx[alias]
	if !ok {
		return nil, persistence.ErrURLNotFound
	}

	return u, nil
}
