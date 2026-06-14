package service

import (
	"context"
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, login, password string) error
	Authenticate(ctx context.Context, login, password string) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
}
