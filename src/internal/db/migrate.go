package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(ctx context.Context, pool *pgxpool.Pool) error {
    sqlContent, err := os.ReadFile("internal/datasource/sql/001_create_users.sql")
    if err != nil {
        return fmt.Errorf("Ошибка чтения файла миграции users: %w", err)
    }
    _, err = pool.Exec(ctx, string(sqlContent))
    if err != nil {
        return fmt.Errorf("миграция users: %w", err)
    }

    sqlContent, err = os.ReadFile("internal/datasource/sql/002_create_games.sql")
    if err != nil {
        return fmt.Errorf("Ошибка чтения файла миграции games: %w", err)
    }
    _, err = pool.Exec(ctx, string(sqlContent))
    if err != nil {
        return fmt.Errorf("миграция games: %w", err)
    }
    return nil
}