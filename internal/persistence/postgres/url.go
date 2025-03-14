package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodeyeen/shortify/internal/domain"
	"github.com/kodeyeen/shortify/internal/persistence"
)

type URLRepository struct {
	dbpool *pgxpool.Pool
}

func NewURLRepository(dbpool *pgxpool.Pool) *URLRepository {
	return &URLRepository{
		dbpool: dbpool,
	}
}

func (r *URLRepository) Close() {
	r.dbpool.Close()
}

func (r *URLRepository) Add(ctx context.Context, u *domain.URL) (int64, error) {
	query := `INSERT INTO urls (original, alias) VALUES (@original, @alias) RETURNING id`
	args := pgx.NamedArgs{
		"original": u.Original,
		"alias":    u.Alias,
	}

	var insertID int64

	err := r.dbpool.QueryRow(ctx, query, args).Scan(&insertID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			switch pgErr.ConstraintName {
			case "urls_original_key":
				return 0, persistence.ErrURLAlreadyExists
			case "urls_alias_key":
				return 0, persistence.ErrDuplicateAlias
			}
		}

		return 0, fmt.Errorf("failed to add url: %w", err)
	}

	return insertID, nil
}

func (r *URLRepository) FindByAlias(ctx context.Context, alias string) (*domain.URL, error) {
	query := `SELECT id, original, alias FROM urls WHERE alias = @alias`
	args := pgx.NamedArgs{
		"alias": alias,
	}

	var u domain.URL

	err := r.dbpool.QueryRow(ctx, query, args).Scan(
		&u.ID,
		&u.Original,
		&u.Alias,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, persistence.ErrURLNotFound
		}

		return nil, fmt.Errorf("failed to find url by alias: %w", err)
	}

	return &u, nil
}
