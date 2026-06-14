package web

import (
	"tic-tac-toe/internal/domain/entity"

	"github.com/google/uuid"
)

type GameDTO struct {
	ID       uuid.UUID                        `json:"id"`
	Board    [entity.Width][entity.Height]int `json:"board"`
	PlayerX  uuid.UUID                        `json:"player_x"`
	PlayerO  uuid.UUID                        `json:"player_o"`
	GameType string                           `json:"game_type"`
	Status   string                           `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
