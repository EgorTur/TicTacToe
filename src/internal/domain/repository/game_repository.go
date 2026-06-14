package repository

import (
	"context"
	"errors"
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
)

type GameRepository interface {
	Save(ctx context.Context, game *entity.Game) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Game, error)
	ListAvailable(ctx context.Context) ([]*entity.Game, error)
	ListCompletedByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Game, error)
	GetTopPlayers(ctx context.Context, limit int) ([]entity.LeaderboardEntry, error)
}

var ErrGameNotFound = errors.New("game not found")
