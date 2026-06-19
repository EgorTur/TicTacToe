package web

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tic-tac-toe/internal/domain/entity"
	"tic-tac-toe/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Вспомогательная функция для добавления userID в контекст запроса
func contextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
    return context.WithValue(ctx, "userID", userID)
}

// --- Тесты CreateGame ---
func TestCreateGame_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    userID := uuid.New()
    game := entity.NewGame(userID, "player")
    game.Status = entity.StatusWaiting

    mockSvc.On("CreateGame", mock.Anything, userID, "player").Return(game, nil)

    reqBody := `{"game_type":"player"}`
    req := httptest.NewRequest(http.MethodPost, "/game", bytes.NewReader([]byte(reqBody)))
    req.Header.Set("Content-Type", "application/json")
    req = req.WithContext(contextWithUserID(req.Context(), userID))
    w := httptest.NewRecorder()

    handler.CreateGame(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    var resp GameDTO
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, game.ID, resp.ID)
    mockSvc.AssertExpectations(t)
}

func TestCreateGame_Unauthorized(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    reqBody := `{"game_type":"player"}`
    req := httptest.NewRequest(http.MethodPost, "/game", bytes.NewReader([]byte(reqBody)))
    req.Header.Set("Content-Type", "application/json")
    // не добавляем userID в контекст
    w := httptest.NewRecorder()

    handler.CreateGame(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    mockSvc.AssertNotCalled(t, "CreateGame")
}

func TestCreateGame_InvalidJSON(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    userID := uuid.New()
    req := httptest.NewRequest(http.MethodPost, "/game", bytes.NewReader([]byte(`not json`)))
    req.Header.Set("Content-Type", "application/json")
    req = req.WithContext(contextWithUserID(req.Context(), userID))
    w := httptest.NewRecorder()

    handler.CreateGame(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    mockSvc.AssertNotCalled(t, "CreateGame")
}

func TestCreateGame_MethodNotAllowed(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    req := httptest.NewRequest(http.MethodGet, "/game", nil)
    w := httptest.NewRecorder()

    handler.CreateGame(w, req)
    assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
    mockSvc.AssertNotCalled(t, "CreateGame")
}

// --- Тесты GetGame ---
func TestGetGame_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    gameID := uuid.New()
    game := entity.NewGame(uuid.New(), "player")
    game.ID = gameID

    mockSvc.On("GetGame", mock.Anything, gameID).Return(game, nil)

    req := httptest.NewRequest(http.MethodGet, "/game/"+gameID.String(), nil)
    w := httptest.NewRecorder()

    handler.GetGame(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    var resp GameDTO
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, gameID, resp.ID)
    mockSvc.AssertExpectations(t)
}

func TestGetGame_NotFound(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    gameID := uuid.New()
    mockSvc.On("GetGame", mock.Anything, gameID).Return(nil, repository.ErrGameNotFound)

    req := httptest.NewRequest(http.MethodGet, "/game/"+gameID.String(), nil)
    w := httptest.NewRecorder()

    handler.GetGame(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
    mockSvc.AssertExpectations(t)
}

// --- Тесты MakeMove ---
func TestMakeMove_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    userID := uuid.New()
    gameID := uuid.New()
    move := map[string]int{"row": 1, "col": 1}
    updatedGame := entity.NewGame(userID, "player")
    updatedGame.ID = gameID
    updatedGame.Status = entity.StatusPlayerOTurn

    mockSvc.On("MakeMove", mock.Anything, gameID, userID, 1, 1).Return(updatedGame, nil)

    body, _ := json.Marshal(move)
    req := httptest.NewRequest(http.MethodPost, "/game/"+gameID.String()+"/move", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req = req.WithContext(contextWithUserID(req.Context(), userID))
    w := httptest.NewRecorder()

    handler.MakeMove(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    mockSvc.AssertExpectations(t)
}

func TestMakeMove_Unauthorized(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    gameID := uuid.New()
    move := map[string]int{"row": 0, "col": 0}
    body, _ := json.Marshal(move)
    req := httptest.NewRequest(http.MethodPost, "/game/"+gameID.String()+"/move", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    handler.MakeMove(w, req)
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    mockSvc.AssertNotCalled(t, "MakeMove")
}

// --- Тесты JoinGame ---
func TestJoinGame_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    joinerID := uuid.New()
    gameID := uuid.New()
    joinedGame := entity.NewGame(uuid.New(), "player")
    joinedGame.ID = gameID
    joinedGame.PlayerO = joinerID
    joinedGame.Status = entity.StatusPlayerXTurn

    mockSvc.On("JoinGame", mock.Anything, gameID, joinerID).Return(joinedGame, nil)

    req := httptest.NewRequest(http.MethodPost, "/game/"+gameID.String()+"/join", nil)
    req = req.WithContext(contextWithUserID(req.Context(), joinerID))
    w := httptest.NewRecorder()

    handler.JoinGame(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    var resp GameDTO
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Equal(t, joinerID, resp.PlayerO)
    mockSvc.AssertExpectations(t)
}

func TestJoinGame_Unauthorized(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    req := httptest.NewRequest(http.MethodPost, "/game/"+uuid.New().String()+"/join", nil)
    w := httptest.NewRecorder()

    handler.JoinGame(w, req)
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    mockSvc.AssertNotCalled(t, "JoinGame")
}

// --- Тесты ListAvailableGames ---
func TestListAvailableGames_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    games := []*entity.Game{
        entity.NewGame(uuid.New(), "player"),
        entity.NewGame(uuid.New(), "player"),
    }
    mockSvc.On("ListAvailableGames", mock.Anything).Return(games, nil)

    req := httptest.NewRequest(http.MethodGet, "/games/list", nil)
    req = req.WithContext(contextWithUserID(req.Context(), uuid.New()))
    w := httptest.NewRecorder()

    handler.ListAvailableGames(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    var dtos []GameDTO
    err := json.Unmarshal(w.Body.Bytes(), &dtos)
    assert.NoError(t, err)
    assert.Len(t, dtos, 2)
    mockSvc.AssertExpectations(t)
}

func TestListAvailableGames_Unauthorized(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    req := httptest.NewRequest(http.MethodGet, "/games/list", nil)
    w := httptest.NewRecorder()

    handler.ListAvailableGames(w, req)
    assert.Equal(t, http.StatusUnauthorized, w.Code)
    mockSvc.AssertNotCalled(t, "ListAvailableGames")
}

// --- Тесты GetCompletedGames ---
func TestGetCompletedGames_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    userID := uuid.New()
    completedGames := []*entity.Game{
        entity.NewGame(userID, "player"),
    }
    mockSvc.On("GetCompletedGames", mock.Anything, userID).Return(completedGames, nil)

    req := httptest.NewRequest(http.MethodGet, "/games/history", nil)
    req = req.WithContext(contextWithUserID(req.Context(), userID))
    w := httptest.NewRecorder()

    handler.GetCompletedGames(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    mockSvc.AssertExpectations(t)
}

// --- Тесты GetTopPlayers ---
func TestGetTopPlayers_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    entries := []entity.LeaderboardEntry{
        {UserID: uuid.New(), Login: "player1", WinRatio: 0.8},
        {UserID: uuid.New(), Login: "player2", WinRatio: 0.5},
    }
    mockSvc.On("GetTopPlayers", mock.Anything, 10).Return(entries, nil)

    req := httptest.NewRequest(http.MethodGet, "/leaderboard?n=10", nil)
    w := httptest.NewRecorder()

    handler.GetTopPlayers(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    var resp []entity.LeaderboardEntry
    err := json.Unmarshal(w.Body.Bytes(), &resp)
    assert.NoError(t, err)
    assert.Len(t, resp, 2)
    mockSvc.AssertExpectations(t)
}

func TestGetTopPlayers_DefaultLimit(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    mockSvc.On("GetTopPlayers", mock.Anything, 10).Return([]entity.LeaderboardEntry{}, nil)

    req := httptest.NewRequest(http.MethodGet, "/leaderboard", nil) // без параметра n
    w := httptest.NewRecorder()

    handler.GetTopPlayers(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    mockSvc.AssertExpectations(t)
}