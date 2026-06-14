package datasource

import (
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
)

func ToModel(game *entity.Game) GameModel {
	var px, po *uuid.UUID
	if game.PlayerX != uuid.Nil {
		px = &game.PlayerX
	}
	if game.PlayerO != uuid.Nil {
		po = &game.PlayerO
	}
	return GameModel{
		ID:        game.ID,
		Board:     boardToFlat(game.Board),
		PlayerX:   px,
		PlayerO:   po,
		GameType:  game.GameType,
		Status:    string(game.Status),
		CreatedAt: game.CreatedAt,
	}
}

func ToEntity(model GameModel) *entity.Game {
	var px, po uuid.UUID
	if model.PlayerX != nil {
		px = *model.PlayerX
	}
	if model.PlayerO != nil {
		po = *model.PlayerO
	}
	return &entity.Game{
		ID:        model.ID,
		Board:     flatToBoard(model.Board),
		PlayerX:   px,
		PlayerO:   po,
		GameType:  model.GameType,
		Status:    entity.GameStatus(model.Status),
		CreatedAt: model.CreatedAt,
	}
}

func ToModelUser(user *entity.User) UserModel {
	return UserModel{
		ID:           user.ID,
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	}
}

func ToEntityUser(model UserModel) *entity.User {
	return &entity.User{
		ID:           model.ID,
		Login:        model.Login,
		PasswordHash: model.PasswordHash,
	}
}

func boardToFlat(b entity.Board) []int {
	flat := make([]int, 0, 9)
	for i := 0; i < entity.Width; i++ {
		for j := 0; j < entity.Height; j++ {
			flat = append(flat, b[i][j])
		}
	}
	return flat
}

func flatToBoard(flat []int) entity.Board {
	var b entity.Board
	for i := 0; i < entity.Width; i++ {
		for j := 0; j < entity.Height; j++ {
			if i*entity.Height+j < len(flat) {
				b[i][j] = flat[i*entity.Height+j]
			}
		}
	}
	return b
}
