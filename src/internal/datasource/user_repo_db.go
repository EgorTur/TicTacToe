package datasource

import (
	"context"
	"errors"
	"fmt"
	"tic-tac-toe/internal/domain/entity"
	"tic-tac-toe/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepoDB struct {
	pool *pgxpool.Pool
}

func NewUserRepoDB(pool *pgxpool.Pool) *UserRepoDB {
	return &UserRepoDB{pool: pool}
}

func (r *UserRepoDB) Create(ctx context.Context, user *entity.User) error {
	model := ToModelUser(user)
	_, err := r.pool.Exec(ctx, usersCreateSQL, model.ID, model.Login, model.PasswordHash)
	if err != nil {
		return fmt.Errorf("Создание пользователя: %w", err)
	}
	return nil
}

func (r *UserRepoDB) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	var model UserModel
	err := r.pool.QueryRow(ctx, getByLoginSQL, login).Scan(
		&model.ID,
		&model.Login,
		&model.PasswordHash,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("get Login: %w", err)
	}
	return ToEntityUser(model), nil
}

func (r *UserRepoDB) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var model UserModel
	err := r.pool.QueryRow(ctx, getByIDSQl, id).Scan(&model.ID, &model.Login, &model.PasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return ToEntityUser(model), nil
}
