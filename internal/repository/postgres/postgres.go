package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(dbURL string) (*Repository, error) {
	const op = "repository.postgres.New"
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: add migrations
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS users(
	    	id BIGSERIAL PRIMARY KEY,
	    	external_user_id TEXT NOT NULL UNIQUE,
	    	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS files(
	    	id UUID PRIMARY KEY,
	    	user_id BIGINT NOT NULL REFERENCES users(id),
	    	path TEXT NOT NULL,
	    	size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0 AND size_bytes <= 1073741824),
	    	mime_type TEXT,
	    	storage_key TEXT NOT NULL,
	    	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	    	UNIQUE (user_id, path)
		);
	`)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{pool: pool}, nil
}
