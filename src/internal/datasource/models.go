package datasource

import (
	"time"

	"github.com/google/uuid"
)

type GameModel struct {
	ID       uuid.UUID
	Board    []int
	PlayerX  *uuid.UUID
	PlayerO  *uuid.UUID
	GameType string
	Status   string
	CreatedAt time.Time
}
