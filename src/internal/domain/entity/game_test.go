package entity

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
    userID := uuid.New()
    game := NewGame(userID, "bot")
    assert.Equal(t, userID, game.PlayerX)
    assert.Equal(t, uuid.Nil, game.PlayerO)
    assert.Equal(t, "bot", game.GameType)
    assert.Equal(t, StatusWaiting, game.Status)
    assert.NotEqual(t, uuid.Nil, game.ID)
}

func TestGame_IsFull(t *testing.T) {
    game := NewGame(uuid.New(), "player")
    assert.False(t, game.IsFull())
    for i := 0; i < Width; i++ {
        for j := 0; j < Height; j++ {
            game.Board.Set(i, j, PlayerX)
        }
    }
    assert.True(t, game.IsFull())
}