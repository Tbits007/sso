package postgres

import (
	"context"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"

	"database/sql"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
    db *sql.DB
}

func New(storagePath string) (*Storage, error) {
    const op = "storage.postgres.New"

    db, err := sql.Open("postgres", storagePath)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
    const op = "storage.postgres.SaveUser"

 
    stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    // Выполняем запрос, передав параметры
    res, err := stmt.ExecContext(ctx, email, passHash)
    if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
		}
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
    const op = "storage.postgres.User"

    stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = $1")
    if err != nil {
        return models.User{}, fmt.Errorf("%s: %w", op, err)
    }

    row := stmt.QueryRowContext(ctx, email)

    var user models.User
    err = row.Scan(&user.ID, &user.Email, &user.PassHash)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
        }

        return models.User{}, fmt.Errorf("%s: %w", op, err)
    }

    return user, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
    const op = "storage.postgres.App"

    stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
    if err != nil {
        return models.App{}, fmt.Errorf("%s: %w", op, err)
    }

    row := stmt.QueryRowContext(ctx, id)

    var app models.App
    err = row.Scan(&app.ID, &app.Name, &app.Secret)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
        }

        return models.App{}, fmt.Errorf("%s: %w", op, err)
    }

    return app, nil
}

