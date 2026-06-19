package mocks

import (
	"context"
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type GameRepositoryMock struct {
    mock.Mock
}

func (m *GameRepositoryMock) Save(ctx context.Context, game *entity.Game) error {
    args := m.Called(ctx, game)
    return args.Error(0)
}

func (m *GameRepositoryMock) Get(ctx context.Context, id uuid.UUID) (*entity.Game, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*entity.Game), args.Error(1)
}

func (m *GameRepositoryMock)ListAvailable(ctx context.Context) ([]*entity.Game, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Game), args.Error(1)
}

func (m *GameRepositoryMock)ListCompletedByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Game, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entity.Game), args.Error(1)
}

func (m *GameRepositoryMock)GetTopPlayers(ctx context.Context, limit int) ([]entity.LeaderboardEntry, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]entity.LeaderboardEntry), args.Error(1)
}