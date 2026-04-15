package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sheixpeer/disk-service/internal/repository"
)

const pgUniqueViolation = "23505"

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

func (r *Repository) CreateUser(ctx context.Context, externalUserID string) (int64, error) {
	const op = "repository.postgres.CreateUser"

	// TODO: мб стоит добавить такую ошибку в repository.go
	if externalUserID == "" {
		return 0, fmt.Errorf("%s: externalUserID is empty", op)
	}

	var id int64

	err := r.pool.QueryRow(ctx, `
		INSERT INTO users (external_user_id)
		VALUES ($1)
		RETURNING id
	`, externalUserID).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return 0, fmt.Errorf("%s: %w", op, repository.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repository) GetUserID(ctx context.Context, externalUserID string) (int64, error) {
	const op = "repository.postgres.GetUserID"

	if externalUserID == "" {
		return 0, fmt.Errorf("%s: externalUserID is empty", op)
	}

	var id int64
	err := r.pool.QueryRow(ctx, `
		SELECT id FROM users WHERE external_user_id = $1
	`, externalUserID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("%s: %w", op, repository.ErrUserNotFound)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// EnsureUserID возвращает внутренний идентификатор пользователя (users.id)
// по внешнему идентификатору externalUserID.
//
// Нужна на текущем этапе, пока нет явной авторизации/регистрации: обработчики могут
// сопоставить внешний идентификатор пользователя со стабильным user_id и выполнять
// операции с файлами, завязанные на FK (files.user_id -> users.id).
//
// Функция безопасна при параллельных вызовах: несколько одновременных запросов с одним
// externalUserID получат один и тот же id; в таблицу будет вставлена максимум одна строка.
func (r *Repository) EnsureUserID(ctx context.Context, externalUserID string) (int64, error) {
	const op = "repository.postgres.EnsureUserID"

	if externalUserID == "" {
		return 0, fmt.Errorf("%s: externalUserID is empty", op)
	}

	var id int64
	err := r.pool.QueryRow(ctx, `
		WITH ins AS (
			INSERT INTO users (external_user_id)
			VALUES ($1)
			ON CONFLICT (external_user_id) DO NOTHING
			RETURNING id
		)
		SELECT id FROM ins
		UNION ALL
		SELECT id FROM users WHERE external_user_id = $1
		LIMIT 1
	`, externalUserID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
