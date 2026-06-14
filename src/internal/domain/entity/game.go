package entity

import (
	"time"

	"github.com/google/uuid"
)

// Константы для клеток
const (
	Empty   = 0
	PlayerX = 1
	PlayerO = 2
)

const (
	StatusWaiting     GameStatus = "waiting"
	StatusPlayerXTurn GameStatus = "player_x_turn"
	StatusPlayerOTurn GameStatus = "player_o_turn"
	StatusDraw        GameStatus = "draw"
	StatusWinX        GameStatus = "win_x"
	StatusWinO        GameStatus = "win_o"
)

type GameStatus string

type Game struct {
	ID       uuid.UUID
	Board    Board
	PlayerX  uuid.UUID
	PlayerO  uuid.UUID
	GameType string
	Status   GameStatus
	CreatedAt time.Time
}

func NewGame(playerX uuid.UUID, gameType string) *Game {
	return &Game{
		ID:       uuid.New(),
		Board:    NewBoard(),
		PlayerX:  playerX,
		PlayerO:  uuid.Nil,
		GameType: gameType,
		Status:   StatusWaiting,
	}

}

func (g *Game) IsFull() bool {
	return g.Board.IsFull()
}
