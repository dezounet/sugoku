package sudoku

import (
	"math/rand"
)

// IsSolved checks if the grid is complete and valid
func (grid *Grid) IsSolved() bool {
	for rowID, columns := range grid.Cells {
		for columnID := range columns {
			cell := grid.GetCell(rowID, columnID)
			value := cell.Value
			cell.Value = empty
			if !grid.IsSafe(rowID, columnID, value) {
				cell.Value = value
				return false
			}

			cell.Value = value
		}
	}

	return true
}

// Solve a sudoku grid
func (grid *Grid) Solve() bool {
	row := 0
	column := 0

	// Is there any empty cell?
	if !grid.findEmptyCell(&row, &column) {
		return true
	}

	// Get some entropy: shuffle values order
	values := make([]int, grid.GetSize())
	for i := 0; i < grid.GetSize(); i++ {
		values[i] = i + 1
	}
	rand.Shuffle(len(values), func(i, j int) { values[i], values[j] = values[j], values[i] })

	// Try every possible value
	for _, value := range values {
		if grid.IsSafe(row, column, value) {
			grid.GetCell(row, column).Value = value

			if grid.Solve() {
				return true
			}

			grid.GetCell(row, column).Value = empty
		}
	}

	return false
}

// CountSolution for a grid, trying all possible valid combinations
// counter should be a valid pointer to an int = 0
func (grid *Grid) CountSolution() int {
	counter := 0

	grid._countSolution(&counter)

	return counter
}

func (grid *Grid) _countSolution(counter *int) {
	row := 0
	column := 0

	if !grid.findEmptyCell(&row, &column) {
		*counter++
	}

	for value := 1; value <= grid.GetSize(); value++ {
		if grid.IsSafe(row, column, value) {
			grid.GetCell(row, column).Value = value

			grid._countSolution(counter)

			grid.GetCell(row, column).Value = empty
		}
	}
}
