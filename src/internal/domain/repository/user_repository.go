package repository

import (
	"context"
	"errors"
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}

var ErrUserNotFound = errors.New("user not found")
