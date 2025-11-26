package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/X/internal/customErr"
	"github.com/saleh-ghazimoradi/X/internal/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	dbWrite *sql.DB
	dbRead  *sql.DB
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at`
	args := []any{user.Username, user.Email, user.Password}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := u.dbWrite.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	var user domain.User

	if err := u.dbRead.QueryRowContext(ctx, query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, customErr.ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user domain.User

	if err := u.dbRead.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, customErr.ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func NewUserRepository(dbWrite, dbRead *sql.DB) UserRepository {
	return &userRepository{
		dbWrite: dbWrite,
		dbRead:  dbRead,
	}
}
