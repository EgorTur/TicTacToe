package web

import (
	"context"
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockGameService struct {
    mock.Mock
}

func (m *mockGameService) CreateGame(ctx context.Context, playerX uuid.UUID, gameType string) (*entity.Game, error) {
    args := m.Called(ctx, playerX, gameType)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Game), args.Error(1)
}

func (m *mockGameService) JoinGame(ctx context.Context, gameID, playerO uuid.UUID) (*entity.Game, error) {
    args := m.Called(ctx, gameID, playerO)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Game), args.Error(1)
}

func (m *mockGameService) MakeMove(ctx context.Context, gameID, playerID uuid.UUID, row, col int) (*entity.Game, error) {
    args := m.Called(ctx, gameID, playerID, row, col)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Game), args.Error(1)
}

func (m *mockGameService) GetGame(ctx context.Context, gameID uuid.UUID) (*entity.Game, error) {
    args := m.Called(ctx, gameID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*entity.Game), args.Error(1)
}

func (m *mockGameService) ListAvailableGames(ctx context.Context) ([]*entity.Game, error) {
    args := m.Called(ctx)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*entity.Game), args.Error(1)
}

func (m *mockGameService) GetCompletedGames(ctx context.Context, userID uuid.UUID) ([]*entity.Game, error) {
    args := m.Called(ctx, userID)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*entity.Game), args.Error(1)
}

func (m *mockGameService) GetTopPlayers(ctx context.Context, limit int) ([]entity.LeaderboardEntry, error) {
    args := m.Called(ctx, limit)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]entity.LeaderboardEntry), args.Error(1)
}