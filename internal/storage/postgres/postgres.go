package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	_ "github.com/lib/pq"
	"todo-sso/internal/domain/models"
	"todo-sso/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func (s *Storage) App(ctx context.Context, appID int) (app models.App, err error) {
	const op = "storage.postgres.App"

	query := `SELECT id, name, secret FROM apps WHERE id = $1`

	row := s.db.QueryRowContext(ctx, query, appID)

	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s %w", op, err)
	}
	return app, nil

}

func New(source string) (*Storage, error) {
	// Пытаемся открыть соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// Проверяем, что соединение действительно установлено и работает
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.postgres.saveUser"

	// Prepare the SQL query with RETURNING clause to get the ID of the inserted row
	query := `INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id`

	// Execute the query with context and parameters
	var id int64
	err := s.db.QueryRowContext(ctx, query, email, passHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 - код ошибки уникальности
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.user"

	query := `SELECT id, email, pass_hash FROM users WHERE email = $1`

	row := s.db.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
