package web

import (
	"errors"
	"tic-tac-toe/internal/domain/entity"
)

func ToEntity(dto *GameDTO) (*entity.Game, error) {
	if dto == nil {
		return nil, errors.New("Ошибка клиента")
	}

	board := entity.NewBoard()

	for i := 0; i < entity.Width; i++ {
		for j := 0; j < entity.Height; j++ {
			board[i][j] = dto.Board[i][j]
		}
	}

	return &entity.Game{
		ID:       dto.ID,
		Board:    board,
		PlayerX:  dto.PlayerX,
		PlayerO:  dto.PlayerO,
		GameType: dto.GameType,
		Status:   entity.GameStatus(dto.Status),
	}, nil
}

func ToDTO(game *entity.Game) (*GameDTO, error) {
	if game == nil {
		return nil, errors.New("Ошибка сервера")
	}

	var board [entity.Width][entity.Height]int

	for i := 0; i < entity.Width; i++ {
		for j := 0; j < entity.Height; j++ {
			board[i][j] = game.Board[i][j]
		}
	}

	return &GameDTO{
		ID:       game.ID,
		Board:    board,
		PlayerX: game.PlayerX,
		PlayerO:   game.PlayerO,
		GameType:   game.GameType,
		Status:  string(game.Status),
	}, nil
}
