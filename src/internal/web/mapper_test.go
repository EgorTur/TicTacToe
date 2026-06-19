package web

import (
    "testing"
    "tic-tac-toe/internal/domain/entity"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func TestToDTO_NilGame(t *testing.T) {
    _, err := ToDTO(nil)
    assert.Error(t, err)
}

func TestToDTO_ValidGame(t *testing.T) {
    id := uuid.New()
    game := &entity.Game{
        ID:       id,
        Board:    entity.NewBoard(),
        PlayerX:  uuid.New(),
        PlayerO:  uuid.New(),
        GameType: "player",
        Status:   entity.StatusWaiting,
    }
    dto, err := ToDTO(game)
    assert.NoError(t, err)
    assert.Equal(t, id, dto.ID)
    assert.Equal(t, game.PlayerX, dto.PlayerX)
    assert.Equal(t, string(game.Status), dto.Status)
}

func TestToEntity_NilDTO(t *testing.T) {
    _, err := ToEntity(nil)
    assert.Error(t, err)
}

func TestToEntity_ValidDTO(t *testing.T) {
    dto := &GameDTO{
        ID:       uuid.New(),
        Board:    [3][3]int{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
        PlayerX:  uuid.New(),
        PlayerO:  uuid.New(),
        GameType: "bot",
        Status:   string(entity.StatusWinX),
    }
    game, err := ToEntity(dto)
    assert.NoError(t, err)
    assert.Equal(t, dto.ID, game.ID)
    assert.Equal(t, entity.GameStatus(dto.Status), game.Status)
    // Проверим, что доска правильно скопировалась
    for i := 0; i < entity.Width; i++ {
        for j := 0; j < entity.Height; j++ {
            assert.Equal(t, dto.Board[i][j], game.Board[i][j])
        }
    }
}