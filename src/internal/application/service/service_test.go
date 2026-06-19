package service

import (
	"context"
	"testing"
	"tic-tac-toe/internal/domain/entity"
	"tic-tac-toe/internal/domain/repository/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateGame_PlayerVsPlayer(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    userID := uuid.New()

    mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.Game")).Return(nil)

    game, err := svc.CreateGame(ctx, userID, "player")
    assert.NoError(t, err)
    assert.Equal(t, entity.StatusWaiting, game.Status)
    assert.Equal(t, userID, game.PlayerX)
    mockRepo.AssertExpectations(t)
}

func TestCreateGame_Bot(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    userID := uuid.New()

    mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.Game")).Return(nil)

    game, err := svc.CreateGame(ctx, userID, "bot")
    assert.NoError(t, err)
    assert.Equal(t, entity.StatusPlayerXTurn, game.Status)
    assert.Equal(t, "bot", game.GameType)
    mockRepo.AssertExpectations(t)
}

func TestJoinGame_Success(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    creator := uuid.New()
    joiner := uuid.New()

    game := entity.NewGame(creator, "player")
    game.Status = entity.StatusWaiting

    mockRepo.On("Get", ctx, game.ID).Return(game, nil)
    mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.Game")).Return(nil)

    joined, err := svc.JoinGame(ctx, game.ID, joiner)
    assert.NoError(t, err)
    assert.Equal(t, joiner, joined.PlayerO)
    assert.Equal(t, entity.StatusPlayerXTurn, joined.Status)
    mockRepo.AssertExpectations(t)
}

func TestJoinGame_NotWaiting(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    game := entity.NewGame(uuid.New(), "player")
    game.Status = entity.StatusPlayerXTurn // уже начата
    mockRepo.On("Get", ctx, game.ID).Return(game, nil)

    _, err := svc.JoinGame(ctx, game.ID, uuid.New())
    assert.ErrorContains(t, err, "session is busy")
    mockRepo.AssertNotCalled(t, "Save") // не должно сохраняться
}

func TestJoinGame_SamePlayer(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    creator := uuid.New()
    game := entity.NewGame(creator, "player")
    game.Status = entity.StatusWaiting
    mockRepo.On("Get", ctx, game.ID).Return(game, nil)

    _, err := svc.JoinGame(ctx, game.ID, creator)
    assert.ErrorContains(t, err, "I can't join my session")
}

func TestMakeMove_ValidMove(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    playerX := uuid.New()
    playerO := uuid.New()

    game := entity.NewGame(playerX, "player")
    game.PlayerO = playerO
    game.Status = entity.StatusPlayerXTurn // ход X

    mockRepo.On("Get", ctx, game.ID).Return(game, nil)
    mockRepo.On("Save", ctx, mock.AnythingOfType("*entity.Game")).Return(nil)

    result, err := svc.MakeMove(ctx, game.ID, playerX, 0, 0)
    assert.NoError(t, err)
    assert.Equal(t, entity.PlayerX, result.Board.Get(0, 0))
    assert.Equal(t, entity.StatusPlayerOTurn, result.Status) // ход перешёл к O
    mockRepo.AssertExpectations(t)
}

func TestMakeMove_WrongPlayer(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    playerX := uuid.New()
    playerO := uuid.New()

    game := entity.NewGame(playerX, "player")
    game.PlayerO = playerO
    game.Status = entity.StatusPlayerXTurn
    mockRepo.On("Get", ctx, game.ID).Return(game, nil)

    _, err := svc.MakeMove(ctx, game.ID, playerO, 0, 0) // ходит O, хотя сейчас ход X
    assert.ErrorContains(t, err, "not your move")
}

func TestMakeMove_OutOfBounds(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    playerX := uuid.New()
    game := entity.NewGame(playerX, "player")
    game.Status = entity.StatusPlayerXTurn
    mockRepo.On("Get", ctx, game.ID).Return(game, nil)

    _, err := svc.MakeMove(ctx, game.ID, playerX, 3, 0) // за пределами
    assert.ErrorContains(t, err, "Incorrect move")
}

func TestMakeMove_CellOccupied(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    playerX := uuid.New()
    game := entity.NewGame(playerX, "player")
    game.Status = entity.StatusPlayerXTurn
    game.Board.Set(0, 0, entity.PlayerX) // уже занята
    mockRepo.On("Get", ctx, game.ID).Return(game, nil)

    _, err := svc.MakeMove(ctx, game.ID, playerX, 0, 0)
    assert.ErrorContains(t, err, "cell is occupied")
}

func TestMakeMove_GameAlreadyOver(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    playerX := uuid.New()
    game := entity.NewGame(playerX, "player")
    game.Status = entity.StatusDraw // игра завершена
    mockRepo.On("Get", ctx, game.ID).Return(game, nil)

    _, err := svc.MakeMove(ctx, game.ID, playerX, 0, 0)
    assert.ErrorContains(t, err, "Incorrect game status")
}

func TestGetGame_Success(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    game := entity.NewGame(uuid.New(), "player")
    mockRepo.On("Get", ctx, game.ID).Return(game, nil)

    result, err := svc.GetGame(ctx, game.ID)
    assert.NoError(t, err)
    assert.Equal(t, game.ID, result.ID)
}

func TestListAvailableGames(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    expected := []*entity.Game{entity.NewGame(uuid.New(), "player")}
    mockRepo.On("ListAvailable", ctx).Return(expected, nil)

    games, err := svc.ListAvailableGames(ctx)
    assert.NoError(t, err)
    assert.Equal(t, expected, games)
}

func TestGetCompletedGames(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    userID := uuid.New()
    expected := []*entity.Game{entity.NewGame(userID, "player")}
    mockRepo.On("ListCompletedByUser", ctx, userID).Return(expected, nil)

    games, err := svc.GetCompletedGames(ctx, userID)
    assert.NoError(t, err)
    assert.Equal(t, expected, games)
}

func TestGetTopPlayers(t *testing.T) {
    mockRepo := new(mocks.GameRepositoryMock)
    svc := NewGameService(mockRepo)
    ctx := context.Background()
    entries := []entity.LeaderboardEntry{{UserID: uuid.New(), Login: "a", WinRatio: 0.8}}
    mockRepo.On("GetTopPlayers", ctx, 5).Return(entries, nil)

    result, err := svc.GetTopPlayers(ctx, 5)
    assert.NoError(t, err)
    assert.Equal(t, entries, result)
}