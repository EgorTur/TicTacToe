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

type gameRepoDB struct {
	pool *pgxpool.Pool
}

func NewGameRepoDB(pool *pgxpool.Pool) repository.GameRepository {
	return &gameRepoDB{pool: pool}
}

func (r *gameRepoDB) Save(ctx context.Context, game *entity.Game) error {
	model := ToModel(game)
	_, err := r.pool.Exec(ctx, createGameSQL,
		model.ID, model.Board, model.PlayerX, model.PlayerO, model.GameType, model.Status)
	if err != nil {
		return fmt.Errorf("save game: %w", err)
	}
	return nil
}

func (r *gameRepoDB) Get(ctx context.Context, id uuid.UUID) (*entity.Game, error) {
	var model GameModel
	err := r.pool.QueryRow(ctx, getGameByIDSQL, id).Scan(
		&model.ID,
		&model.Board,
		&model.PlayerX,
		&model.PlayerO,
		&model.GameType,
		&model.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrGameNotFound
		}
		return nil, fmt.Errorf("get game: %w", err)
	}
	return ToEntity(model), nil
}

func (r *gameRepoDB) ListAvailable(ctx context.Context) ([]*entity.Game, error) {
	rows, err := r.pool.Query(ctx, listAvailableGamesSQL)
	if err != nil {
		return nil, fmt.Errorf("list available games: %w", err)
	}
	defer rows.Close()

	var games []*entity.Game
	for rows.Next() {
		var model GameModel
		err := rows.Scan(
			&model.ID,
			&model.Board,
			&model.PlayerX,
			&model.PlayerO,
			&model.GameType,
			&model.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scan game: %w", err)
		}
		games = append(games, ToEntity(model))

	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return games, nil
}

func (r *gameRepoDB) ListCompletedByUser(ctx context.Context, userID uuid.UUID) ([]*entity.Game, error) {
	rows, err := r.pool.Query(ctx, listCompletedByUserSQL, userID)
	if err != nil {
		return nil, fmt.Errorf("list completed games by user: %w", err)
	}
	defer rows.Close()
	var games []*entity.Game
	for rows.Next() {
		var model GameModel
		err := rows.Scan(
			&model.ID,
			&model.Board,
			&model.PlayerX,
			&model.PlayerO,
			&model.GameType,
			&model.Status,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan game: %w", err)
		}
		games = append(games, ToEntity(model))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return games, nil
}

func (r *gameRepoDB) GetTopPlayers(ctx context.Context, limit int) ([]entity.LeaderboardEntry, error) {
	rows, err := r.pool.Query(ctx, leaderboardSQL, limit)
	if err != nil {
		return nil, fmt.Errorf("get top players: %w", err)
	}
	defer rows.Close()

	var entries []entity.LeaderboardEntry
	for rows.Next() {
		var entry entity.LeaderboardEntry
		if err := rows.Scan(&entry.UserID, &entry.Login, &entry.WinRatio); err != nil {
			return nil, fmt.Errorf("scan leaderboard entry: %w", err)
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}