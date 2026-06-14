package datasource

import (
	"embed"
	"fmt"
)

//go:embed sql/*.sql
var sqlFiles embed.FS

var (
	createGameSQL          string
	getGameByIDSQL         string
	updatedGameSQL         string
	usersCreateSQL         string
	getByLoginSQL          string
	listAvailableGamesSQL  string
	listCompletedByUserSQL string
	getByIDSQl             string
	leaderboardSQL         string
)

func init() {
	var err error
	createGameSQL, err = readSQLFile("sql/games_create.sql")
	if err != nil {
		panic(err)
	}

	getGameByIDSQL, err = readSQLFile("sql/games_get_by_id.sql")
	if err != nil {
		panic(err)
	}

	updatedGameSQL, err = readSQLFile("sql/games_update.sql")
	if err != nil {
		panic(err)
	}
	usersCreateSQL, err = readSQLFile("sql/user_create.sql")
	if err != nil {
		panic(err)
	}
	getByLoginSQL, err = readSQLFile("sql/users_get_by_login.sql")
	if err != nil {
		panic(err)
	}
	listAvailableGamesSQL, err = readSQLFile("sql/games_list_available.sql")
	if err != nil {
		panic(err)
	}
	getByIDSQl, err = readSQLFile("sql/get_id.sql")
	if err != nil {
		panic(err)
	}
	listCompletedByUserSQL, err = readSQLFile("sql/games_comleted_by_user.sql")
	if err != nil {
		panic(err)
	}
	leaderboardSQL, err = readSQLFile("sql/leaderboard.sql")
	if err != nil {
		panic(err)
	}
}

func readSQLFile(name string) (string, error) {
	data, err := sqlFiles.ReadFile(name)
	if err != nil {
		return "", fmt.Errorf("read sql file %s: %w", name, err)
	}
	return string(data), nil
}
