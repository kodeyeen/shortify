package inmemory

import (
	"context"
	"sync"

	"github.com/kodeyeen/shortify/internal/domain"
	"github.com/kodeyeen/shortify/internal/persistence"
)

type URLRepository struct {
	items map[int64]*domain.URL

	aliasIdx map[string]*domain.URL

	mu *sync.RWMutex
}

func NewURLRepository() *URLRepository {
	return &URLRepository{
		items: map[int64]*domain.URL{},

		aliasIdx: map[string]*domain.URL{},

		mu: &sync.RWMutex{},
	}
}

func (r *URLRepository) Add(ctx context.Context, u *domain.URL) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.items[u.ID]; ok {
		return 0, persistence.ErrURLAlreadyExists
	}

	id := int64(len(r.items) + 1)

	r.items[id] = u
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
