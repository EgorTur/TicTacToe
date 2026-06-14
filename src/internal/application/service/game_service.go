package service

import (
	"context"
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
)

type GameService interface {
	CreateGame(ctx context.Context, playerX uuid.UUID, gameType string) (*entity.Game, error)
	JoinGame(ctx context.Context, gameID uuid.UUID, playerO uuid.UUID) (*entity.Game, error)
	MakeMove(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID, row, col int) (*entity.Game, error)
	GetGame(ctx context.Context, gameID uuid.UUID) (*entity.Game, error)
	ListAvailableGames(ctx context.Context) ([]*entity.Game, error)
	GetCompletedGames(ctx context.Context, userID uuid.UUID) ([]*entity.Game, error)
	GetTopPlayers(ctx context.Context, limit int) ([]entity.LeaderboardEntry, error)
}
