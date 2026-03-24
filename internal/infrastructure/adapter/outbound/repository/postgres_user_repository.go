package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/apex20/backend/internal/application/port"
	"github.com/apex20/backend/internal/domain/user"
	repositorygen "github.com/apex20/backend/internal/infrastructure/adapter/outbound/repository/gen"
)

type PostgresUserRepository struct {
	queries *repositorygen.Queries
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{queries: repositorygen.New(db)}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, u user.User) error {
	err := r.queries.CreateUser(ctx, repositorygen.CreateUserParams{
		ID:           u.ID,
		Email:        u.Email,
		Name:         u.Name,
		Nick:         sql.NullString{String: u.Nick, Valid: u.Nick != ""},
		PasswordHash: u.PasswordHash,
		IsAdmin:      u.IsAdmin,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return port.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (user.User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, port.ErrNotFound
		}
		return user.User{}, err
	}
	return toUserDomain(row), nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	row, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, port.ErrNotFound
		}
		return user.User{}, err
	}
	return user.User{
		ID:        row.ID,
		Email:     row.Email,
		Name:      row.Name,
		Nick:      row.Nick.String,
		IsAdmin:   row.IsAdmin,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, id uuid.UUID, name string, nick *string) (user.User, error) {
	var nullNick sql.NullString
	if nick != nil {
		nullNick = sql.NullString{String: *nick, Valid: true}
	}
	row, err := r.queries.UpdateUser(ctx, repositorygen.UpdateUserParams{
		ID:   id,
		Name: name,
		Nick: nullNick,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user.User{}, port.ErrNotFound
		}
		return user.User{}, err
	}
	return user.User{
		ID:        row.ID,
		Email:     row.Email,
		Name:      row.Name,
		Nick:      row.Nick.String,
		IsAdmin:   row.IsAdmin,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	n, err := r.queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	if n == 0 {
		return port.ErrNotFound
	}
	return nil
}

func toUserDomain(row repositorygen.GetUserByEmailRow) user.User {
	return user.User{
		ID:           row.ID,
		Email:        row.Email,
		Name:         row.Name,
		Nick:         row.Nick.String,
		PasswordHash: row.PasswordHash,
		IsAdmin:      row.IsAdmin,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}
