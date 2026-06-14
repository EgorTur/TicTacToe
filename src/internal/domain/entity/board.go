 package entity

const (
	Width  = 3
	Height = 3
)

type Board [Width][Height]int

func NewBoard() Board {
	return [Width][Height]int{}
}

// Установка значения в клетку
func (b *Board) Set(row, col, value int) {
	b[row][col] = value
}

// Возврат значения из клетки
func (b *Board) Get(row, col int) int {
	return b[row][col]
}

// Проверка пустой клетки
func (b *Board) IsEmpty(row, col int) bool {
	if b[row][col] == 0 {
		return true
	}
	return false
}

// Возврат всех пустых клеток
func (b *Board) GetAllEmpty() [][]int {
	var emptyCell [][]int
	for i := 0; i < Width; i++ {

		for j := 0; j < Height; j++ {
			if b.IsEmpty(i, j) {
				emptyCell = append(emptyCell, []int{i, j})
			}
		}
	}
	return emptyCell
}

func (b *Board) IsFull() bool {
	for i := 0; i < Width; i++ {
		for j := 0; j < Height; j++ {
			if b[i][j] == 0 {
				return false
			}
		}
	}
	return true
}

