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

// Вспомогательная функция для добавления userID в контекст
func contextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
    return context.WithValue(ctx, "userID", userID)
}

// Вспомогательная функция для выполнения запроса через mux, чтобы работал PathValue
func serveHandler(mux *http.ServeMux, method, path string, body []byte, userID *uuid.UUID) *httptest.ResponseRecorder {
    var req *http.Request
    if body != nil {
        req = httptest.NewRequest(method, path, bytes.NewReader(body))
    } else {
        req = httptest.NewRequest(method, path, nil)
    }
    req.Header.Set("Content-Type", "application/json")
    if userID != nil {
        req = req.WithContext(contextWithUserID(req.Context(), *userID))
    }
    w := httptest.NewRecorder()
    mux.ServeHTTP(w, req)
    return w
}

// --- CreateGame ---
func TestCreateGame_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    userID := uuid.New()
    game := entity.NewGame(userID, "player")
    game.Status = entity.StatusWaiting

    mockSvc.On("CreateGame", mock.Anything, userID, "player").Return(game, nil)

    mux := http.NewServeMux()
    mux.HandleFunc("POST /game", handler.CreateGame)

    reqBody := `{"game_type":"player"}`
    w := serveHandler(mux, "POST", "/game", []byte(reqBody), &userID)

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

    mux := http.NewServeMux()
    mux.HandleFunc("POST /game", handler.CreateGame)

    reqBody := `{"game_type":"player"}`
    w := serveHandler(mux, "POST", "/game", []byte(reqBody), nil)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    mockSvc.AssertNotCalled(t, "CreateGame")
}

// --- GetGame ---
func TestGetGame_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    gameID := uuid.New()
    game := entity.NewGame(uuid.New(), "player")
    game.ID = gameID

    mockSvc.On("GetGame", mock.Anything, gameID).Return(game, nil)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /game/{id}", handler.GetGame)

    w := serveHandler(mux, "GET", "/game/"+gameID.String(), nil, nil)

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

    mux := http.NewServeMux()
    mux.HandleFunc("GET /game/{id}", handler.GetGame)

    w := serveHandler(mux, "GET", "/game/"+gameID.String(), nil, nil)

    assert.Equal(t, http.StatusNotFound, w.Code)
    mockSvc.AssertExpectations(t)
}

// --- MakeMove ---
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

    mux := http.NewServeMux()
    // Используем путь, соответствующий маршруту в сервере: POST /game/{id}/move
    mux.HandleFunc("POST /game/{id}/move", handler.MakeMove)

    body, _ := json.Marshal(move)
    w := serveHandler(mux, "POST", "/game/"+gameID.String()+"/move", body, &userID)

    assert.Equal(t, http.StatusOK, w.Code)
    mockSvc.AssertExpectations(t)
}

func TestMakeMove_Unauthorized(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    gameID := uuid.New()
    move := map[string]int{"row": 0, "col": 0}
    body, _ := json.Marshal(move)

    mux := http.NewServeMux()
    mux.HandleFunc("POST /game/{id}/move", handler.MakeMove)
    w := serveHandler(mux, "POST", "/game/"+gameID.String()+"/move", body, nil)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    mockSvc.AssertNotCalled(t, "MakeMove")
}

// --- JoinGame ---
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

    mux := http.NewServeMux()
    mux.HandleFunc("POST /game/{id}/join", handler.JoinGame)

    w := serveHandler(mux, "POST", "/game/"+gameID.String()+"/join", nil, &joinerID)

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

    mux := http.NewServeMux()
    mux.HandleFunc("POST /game/{id}/join", handler.JoinGame)
    w := serveHandler(mux, "POST", "/game/"+uuid.New().String()+"/join", nil, nil)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    mockSvc.AssertNotCalled(t, "JoinGame")
}

// --- ListAvailableGames ---
func TestListAvailableGames_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    games := []*entity.Game{
        entity.NewGame(uuid.New(), "player"),
        entity.NewGame(uuid.New(), "player"),
    }
    mockSvc.On("ListAvailableGames", mock.Anything).Return(games, nil)

    userID := uuid.New()
    mux := http.NewServeMux()
    mux.HandleFunc("GET /game/list", handler.ListAvailableGames)

    w := serveHandler(mux, "GET", "/game/list", nil, &userID)

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

    mux := http.NewServeMux()
    mux.HandleFunc("GET /game/list", handler.ListAvailableGames)
    w := serveHandler(mux, "GET", "/game/list", nil, nil)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    mockSvc.AssertNotCalled(t, "ListAvailableGames")
}

// --- GetCompletedGames ---
func TestGetCompletedGames_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    userID := uuid.New()
    completed := []*entity.Game{entity.NewGame(userID, "player")}
    mockSvc.On("GetCompletedGames", mock.Anything, userID).Return(completed, nil)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /games/history", handler.GetCompletedGames)

    w := serveHandler(mux, "GET", "/games/history", nil, &userID)

    assert.Equal(t, http.StatusOK, w.Code)
    mockSvc.AssertExpectations(t)
}

// --- GetTopPlayers ---
func TestGetTopPlayers_Success(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    entries := []entity.LeaderboardEntry{
        {UserID: uuid.New(), Login: "p1", WinRatio: 0.8},
    }
    mockSvc.On("GetTopPlayers", mock.Anything, 10).Return(entries, nil)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /leaderboard", handler.GetTopPlayers)

    w := serveHandler(mux, "GET", "/leaderboard?n=10", nil, nil)

    assert.Equal(t, http.StatusOK, w.Code)
    var resp []entity.LeaderboardEntry
    json.Unmarshal(w.Body.Bytes(), &resp)
    assert.Len(t, resp, 1)
    mockSvc.AssertExpectations(t)
}

func TestGetTopPlayers_DefaultLimit(t *testing.T) {
    mockSvc := new(mockGameService)
    handler := NewGameHandler(mockSvc)

    mockSvc.On("GetTopPlayers", mock.Anything, 10).Return([]entity.LeaderboardEntry{}, nil)

    mux := http.NewServeMux()
    mux.HandleFunc("GET /leaderboard", handler.GetTopPlayers)

    w := serveHandler(mux, "GET", "/leaderboard", nil, nil)

    assert.Equal(t, http.StatusOK, w.Code)
    mockSvc.AssertExpectations(t)
}