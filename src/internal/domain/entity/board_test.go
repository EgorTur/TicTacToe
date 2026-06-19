package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBoard(t *testing.T) {
    b := NewBoard()
    for i := 0; i < Width; i++ {
        for j := 0; j < Height; j++ {
            assert.Equal(t, Empty, b[i][j])
        }
    }
}

func TestBoard_SetGet(t *testing.T) {
    b := NewBoard()
    b.Set(0, 2, PlayerX)
    assert.Equal(t, PlayerX, b.Get(0, 2))
}

func TestBoard_IsEmpty(t *testing.T) {
    b := NewBoard()
    assert.True(t, b.IsEmpty(1, 1))
    b.Set(1, 1, PlayerO)
    assert.False(t, b.IsEmpty(1, 1))
}

func TestBoard_GetAllEmpty(t *testing.T) {
    b := NewBoard()
    // Заполним все клетки, кроме одной
    for i := 0; i < Width; i++ {
        for j := 0; j < Height; j++ {
            b.Set(i, j, PlayerX)
        }
    }
    b.Set(2, 2, Empty)
    empty := b.GetAllEmpty()
    assert.Len(t, empty, 1)
    assert.Equal(t, []int{2, 2}, empty[0])
}

func TestBoard_IsFull(t *testing.T) {
    b := NewBoard()
    assert.False(t, b.IsFull())
    for i := 0; i < Width; i++ {
        for j := 0; j < Height; j++ {
            b.Set(i, j, PlayerX)
        }
    }
    assert.True(t, b.IsFull())
}