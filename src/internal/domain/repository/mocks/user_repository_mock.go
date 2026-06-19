package mocks

import (
	"context"
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
    mock.Mock
}

func (m *UserRepositoryMock) Create(ctx context.Context, user *entity.User) error {
    args := m.Called(ctx, user)
    return args.Error(0)
}

func (m *UserRepositoryMock) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
    args := m.Called(ctx, login)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.User), args.Error(1)
}