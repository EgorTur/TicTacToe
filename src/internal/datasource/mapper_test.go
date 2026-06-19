package datasource

import (
	"testing"
	"tic-tac-toe/internal/domain/entity"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBoardToFlat(t *testing.T) {
    board := entity.NewBoard()
    board.Set(0, 0, 1)
    board.Set(2, 2, 2)
    flat := boardToFlat(board)
    assert.Len(t, flat, 9)
    assert.Equal(t, 1, flat[0])
    assert.Equal(t, 0, flat[1])
    assert.Equal(t, 2, flat[8])
}

func TestFlatToBoard(t *testing.T) {
    flat := []int{0, 1, 0, 0, 0, 0, 0, 0, 2}
    board := flatToBoard(flat)
    assert.Equal(t, 1, board[0][1])
    assert.Equal(t, 2, board[2][2])
}

func TestToModelAndToEntity(t *testing.T) {
    userX := uuid.New()
    userO := uuid.New()
    now := time.Now().Truncate(time.Second)
    game := &entity.Game{
        ID:        uuid.New(),
        Board:     entity.NewBoard(),
        PlayerX:   userX,
        PlayerO:   userO,
        GameType:  "player",
        Status:    entity.StatusPlayerXTurn,
        CreatedAt: now,
    }
    game.Board.Set(1, 1, 1)

    model := ToModel(game)
    assert.Equal(t, game.ID, model.ID)
    assert.Equal(t, userX, *model.PlayerX)
    assert.Equal(t, userO, *model.PlayerO)
    assert.Equal(t, string(game.Status), model.Status)
    assert.WithinDuration(t, now, model.CreatedAt, 0)

    // Обратно в entity
    gameBack := ToEntity(model)
    assert.Equal(t, game.ID, gameBack.ID)
    assert.Equal(t, game.PlayerX, gameBack.PlayerX)
    assert.Equal(t, game.PlayerO, gameBack.PlayerO)
    assert.Equal(t, game.Status, gameBack.Status)
    assert.Equal(t, game.Board, gameBack.Board)
    assert.WithinDuration(t, game.CreatedAt, gameBack.CreatedAt, 0)
}

func TestToModelUser_ToEntityUser(t *testing.T) {
    user := &entity.User{
        ID:           uuid.New(),
        Login:        "testuser",
        PasswordHash: "hash",
    }
    model := ToModelUser(user)
    assert.Equal(t, user.ID, model.ID)
    assert.Equal(t, user.Login, model.Login)
    assert.Equal(t, user.PasswordHash, model.PasswordHash)

    userBack := ToEntityUser(model)
    assert.Equal(t, user.ID, userBack.ID)
    assert.Equal(t, user.Login, userBack.Login)
    assert.Equal(t, user.PasswordHash, userBack.PasswordHash)
}