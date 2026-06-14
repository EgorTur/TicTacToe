package service

import (
	"context"
	"errors"
	"fmt"
	"tic-tac-toe/internal/domain/entity"
	"tic-tac-toe/internal/domain/repository"

	"github.com/google/uuid"
)

const (
	scoreBotWin    = 10
	scorePlayerWin = -10
	scoreDraw      = 0
)

type gameService struct {
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) *gameService {
	return &gameService{repo: repo}
}

func (s *gameService) CreateGame(ctx context.Context, playerX uuid.UUID, gameType string) (*entity.Game, error) {
	if gameType != "player" && gameType != "bot" {
		gameType = "player"
	}
	game := entity.NewGame(playerX, gameType)
	if gameType == "bot" {
		game.Status = entity.StatusPlayerXTurn
	}
	if err := s.SaveGame(ctx, game); err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) JoinGame(ctx context.Context, gameID uuid.UUID, playerO uuid.UUID) (*entity.Game, error) {
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return nil, err
	}
	if game.Status != entity.StatusWaiting {
		return nil, errors.New("session is busy")
	}
	if playerO == game.PlayerX {
		return nil, errors.New("I can't join my session")
	}
	game.PlayerO = playerO
	game.Status = entity.StatusPlayerXTurn
	if err := s.SaveGame(ctx, game); err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) MakeMove(ctx context.Context, gameID uuid.UUID, playerID uuid.UUID, row, col int) (*entity.Game, error) {
	game, err := s.GetGame(ctx, gameID)
	if err != nil {
		return nil, err
	}

	if row < 0 || row >= entity.Width || col < 0 ||
		col >= entity.Height {
		return nil, errors.New("Incorrect move")
	}

	if game.Status != entity.StatusPlayerOTurn && game.Status != entity.StatusPlayerXTurn {
		return nil, errors.New("Incorrect game status")
	}

	if game.Status == entity.StatusPlayerXTurn && playerID != game.PlayerX ||
		game.Status == entity.StatusPlayerOTurn && playerID != game.PlayerO {
		return nil, errors.New("not your move")
	}

	if !game.Board.IsEmpty(row, col) {
		return nil, errors.New("cell is occupied")
	}

	playerSign := entity.PlayerX
	if game.Status == entity.StatusPlayerOTurn {
		playerSign = entity.PlayerO
	}

	game.Board.Set(row, col, playerSign)

	s.updateStatus(game)

	if game.GameType == "bot" && game.Status == entity.StatusPlayerOTurn {
		botRow, botCol, err := s.GetNextMove(*game)
		if err != nil {
			return nil, fmt.Errorf("bot move failed: %w", err)
		}
		game.Board.Set(botRow, botCol, entity.PlayerO)
		s.updateStatus(game)
	}

	if err := s.SaveGame(ctx, game); err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) updateStatus(game *entity.Game) {
	if winner, ok := checkWinner(game.Board); ok {
		if winner == entity.PlayerX {
			game.Status = entity.StatusWinX
		} else {
			game.Status = entity.StatusWinO
		}
	} else if game.Board.IsFull() {
		game.Status = entity.StatusDraw
	} else {
		if game.Status == entity.StatusPlayerXTurn {
			game.Status = entity.StatusPlayerOTurn
		} else {
			game.Status = entity.StatusPlayerXTurn
		}
	}
}

func copyBoard(original entity.Board) entity.Board {
	newBoard := entity.NewBoard()
	for i := 0; i < entity.Width; i++ {
		for j := 0; j < entity.Height; j++ {
			newBoard[i][j] = original[i][j]
		}
	}
	return newBoard
}

func minimax(board entity.Board, depth int, isMaxmizing bool) int {

	sig, _ := checkWinner(board)

	if sig == 2 {
		return scoreBotWin
	}

	if sig == 1 {
		return scorePlayerWin
	}

	if board.IsFull() {
		return scoreDraw
	}

	if depth == 0 {
		return scoreDraw
	}

	emptyCells := board.GetAllEmpty()

	if isMaxmizing {
		// Ход бота
		best := -1000

		for _, cell := range emptyCells {
			row := cell[0]
			col := cell[1]

			newBoard := copyBoard(board)
			newBoard.Set(row, col, 2) // 2 = бот
			score := minimax(newBoard, depth-1, false)

			if score > best {
				best = score
			}
		}
		return best

	} else {
		// Ход игрока
		best := 1000

		for _, cell := range emptyCells {
			row := cell[0]
			col := cell[1]

			newBoard := copyBoard(board)
			newBoard.Set(row, col, 1) // 1 = игрок
			score := minimax(newBoard, depth-1, true)

			if score < best {
				best = score
			}
		}
		return best
	}
}

func (s *gameService) GetNextMove(game entity.Game) (row, col int, err error) {
	emptyCells := game.Board.GetAllEmpty()
	if len(emptyCells) == 0 {
		return -1, -1, errors.New("Нету пустых клеток")
	}

	bestScore := -1000
	bestRow, bestCol := -1, -1
	depth := len(emptyCells)

	for _, cell := range emptyCells {
		r, c := cell[0], cell[1]
		newBoard := copyBoard(game.Board)
		newBoard.Set(r, c, entity.PlayerO)
		score := minimax(newBoard, depth-1, false)

		if score > bestScore {
			bestScore = score
			bestRow, bestCol = r, c
		}
	}
	return bestRow, bestCol, nil
}

// CheckWinner - проверяет наличие победителя на доске
func checkWinner(board entity.Board) (int, bool) {
	for i := 0; i < entity.Width; i++ {
		if board[i][0] != entity.Empty && board[i][0] == board[i][1] && board[i][1] == board[i][2] {
			return board[i][0], true
		}
	}
	for j := 0; j < entity.Height; j++ {
		if board[0][j] != entity.Empty && board[0][j] == board[1][j] && board[1][j] == board[2][j] {
			return board[0][j], true
		}
	}
	if board[0][0] != entity.Empty && board[0][0] == board[1][1] && board[1][1] == board[2][2] {
		return board[0][0], true
	}
	if board[0][2] != entity.Empty && board[0][2] == board[1][1] && board[1][1] == board[2][0] {
		return board[0][2], true
	}
	return entity.Empty, false
}

func (s *gameService) GetGame(ctx context.Context, id uuid.UUID) (*entity.Game, error) {
	return s.repo.Get(ctx, id)
}

func (s *gameService) SaveGame(ctx context.Context, game *entity.Game) error {
	return s.repo.Save(ctx, game)
}

func (s *gameService) ListAvailableGames(ctx context.Context) ([]*entity.Game, error) {
	return s.repo.ListAvailable(ctx)
}

func (s *gameService) GetCompletedGames(ctx context.Context, userID uuid.UUID) ([]*entity.Game, error) {
	return s.repo.ListCompletedByUser(ctx, userID)
}

func (s *gameService) GetTopPlayers(ctx context.Context, limit int) ([]entity.LeaderboardEntry, error) {
	return s.repo.GetTopPlayers(ctx, limit)
}
