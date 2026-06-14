package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"tic-tac-toe/internal/application/service"
	"tic-tac-toe/internal/domain/entity"
	"tic-tac-toe/internal/domain/repository"

	"github.com/google/uuid"
)

type GameHandler struct {
	service service.GameService
}

func NewGameHandler(service service.GameService) *GameHandler {
	return &GameHandler{
		service: service,
	}
}

func (g *GameHandler) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		GameType string `json:"game_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "incorrect JSON")
		return
	}

	userID, ok := r.Context().Value("userID").(uuid.UUID)

	if !ok {
		writeError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	game, err := g.service.CreateGame(r.Context(), userID, req.GameType)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	CreateRespons(game, w)
}

func (g *GameHandler) GetGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid game id")
		return
	}

	currentGame, err := g.service.GetGame(r.Context(), id)

	if err != nil {
		if errors.Is(err, repository.ErrGameNotFound) {

			writeError(w, http.StatusNotFound, "game not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	CreateRespons(currentGame, w)

}

func (g *GameHandler) MakeMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	gameID, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid game id")
		return
	}

	userID, ok := r.Context().Value("userID").(uuid.UUID)

	if !ok {
		writeError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	var move struct {
		Row int `json:"row"`
		Col int `json:"col"`
	}

	if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	game, err := g.service.MakeMove(r.Context(), gameID, userID, move.Row, move.Col)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	CreateRespons(game, w)
}

func (g *GameHandler) JoinGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	gameID, err := uuid.Parse(r.PathValue("id"))

	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid game id")
		return
	}

	userID, ok := r.Context().Value("userID").(uuid.UUID)

	if !ok {
		writeError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	game, err := g.service.JoinGame(r.Context(), gameID, userID)

	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	CreateRespons(game, w)
}

func (g *GameHandler) ListAvailableGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if _, ok := r.Context().Value("userID").(uuid.UUID); !ok {
		writeError(w, http.StatusUnauthorized, "Not authorized")
		return
	}

	games, err := g.service.ListAvailableGames(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	CreateGamesListResponse(games, w)
}

func CreateGamesListResponse(games []*entity.Game, w http.ResponseWriter) {
	dtos := make([]GameDTO, 0, len(games))
	for _, game := range games {
		dto, err := ToDTO(game)
		if err != nil {
			continue
		}
		dtos = append(dtos, *dto)

	}
	if dtos == nil {
		dtos = []GameDTO{}
	}
	writeJSON(w, http.StatusOK, dtos)

}

func (g *GameHandler) GetCompletedGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		writeError(w, http.StatusUnauthorized, "not authorized")
		return
	}
	games, err := g.service.GetCompletedGames(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	CreateGamesListResponse(games, w)
}

func (g *GameHandler) GetTopPlayers(w http.ResponseWriter, r *http.Request) {
	nStr := r.URL.Query().Get("n")
	n := 10
	if nStr != "" {
		if val, err := strconv.Atoi(nStr); err == nil && val > 0 {
			n = val
		}
	}
	entries, err := g.service.GetTopPlayers(r.Context(), n)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func CreateRespons(game *entity.Game, w http.ResponseWriter) {
	responsDTO, err := ToDTO(game)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, responsDTO)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}
