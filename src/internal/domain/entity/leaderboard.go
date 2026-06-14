package entity

import "github.com/google/uuid"

type LeaderboardEntry struct {
	UserID   uuid.UUID `json:"user_id"`
	Login    string    `json:"login"`
	WinRatio float64   `json:"win_ratio"`
}
