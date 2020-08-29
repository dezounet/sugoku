package sudoku

import "github.com/google/uuid"

// isEmpty checks if a cell has no valid value yet
func (cell *Cell) isEmpty() bool {
	return cell.Value == empty
}

// GetSize of the grid
func (grid *Grid) GetSize() int {
	return grid.BlockSize * grid.BlockSize
}

// GetCell from a row and column IDs
func (grid *Grid) GetCell(row int, column int) *Cell {
	return &grid.Cells[row][column]
}

// CreateEmptyGrid of choosen size
func CreateEmptyGrid(blockSize int) *Grid {
	grid := Grid{
		BlockSize: blockSize,
	}

	// Allocate cells
	rows := make([][]Cell, grid.GetSize())
	for rowID := range rows {
		columns := make([]Cell, grid.GetSize())
		rows[rowID] = columns
	}
	grid.Cells = rows

	// Initialize UUID & cells
	grid.Reset()

	return &grid
}

// Reset a grid UUID and its cells
func (grid *Grid) Reset() {
	grid.UUID = uuid.New().String()

	for i := range grid.Cells {
		for j := range grid.Cells[i] {
			cell := Cell{
				Coordinates: Coordinates{
					Row:    i,
					Column: j,
				},
			}

			grid.Cells[i][j] = cell
		}
	}
}

// findEmptyCell in a grid
func (grid *Grid) findEmptyCell(row *int, column *int) bool {
	for rowID, columns := range grid.Cells {
		for columnID, cell := range columns {
			if cell.isEmpty() {
				*row = rowID
				*column = columnID

				return true
			}
		}
	}

	return false
}

// CountEmptyCell in a grid
func (grid *Grid) CountEmptyCell() int {
	counter := 0

	for rowID, columns := range grid.Cells {
		for columnID := range columns {
			if grid.GetCell(rowID, columnID).isEmpty() {
				counter++
			}
		}
	}

	return counter
}
