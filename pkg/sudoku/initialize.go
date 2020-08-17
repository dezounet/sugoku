package sudoku

import (
	"math/rand"
)

// DifficultyToEmptyCellCount to get a number of cell to empty
func (grid *Grid) DifficultyToEmptyCellCount(difficulty Difficulty) int {
	emptyCellCount := 0

	totalCellCount := grid.GetSize() * grid.GetSize()
	switch difficulty {
	case NIGHTMARE:
		emptyCellCount = 70 * totalCellCount / 100
	case HARD:
		emptyCellCount = 67 * totalCellCount / 100
	case MEDIUM:
		emptyCellCount = 64 * totalCellCount / 100
	case EASY:
		emptyCellCount = 61 * totalCellCount / 100
	default:
		emptyCellCount = 61 * totalCellCount / 100
	}

	return emptyCellCount
}

// Initialize a grid with a given difficulty
func (grid *Grid) Initialize(difficulty Difficulty) {
	// Select number of empty cell from expected difficulty
	expectedEmptyCellCount := grid.DifficultyToEmptyCellCount(difficulty)

	// Erase everything in grid and solve it
	grid.Reset()
	grid.Solve()

	// Delete known cell values while only 1 solution exists,
	// keep going until more than one solution appears or
	// retry on a new grid
	for !grid.RemoveCellValue(expectedEmptyCellCount) {
		grid.Reset()
		grid.Solve()
	}

	for rowID := range grid.Cells {
		for columnID := range grid.Cells[rowID] {
			cell := grid.GetCell(rowID, columnID)
			if !cell.isEmpty() {
				cell.Frozen = true
			}
		}
	}
}

// RemoveCellValue cell values by deleting them, until we are barely left with a
// 1 solution-only sudoku. Remaining cells are flagged as frozen.
func (grid *Grid) RemoveCellValue(emptyCellCount int) bool {
	if emptyCellCount <= grid.GetSize()*grid.GetSize() && grid.Solve() {
		cells := []Cell{}
		for _, columns := range grid.Cells {
			for _, cell := range columns {
				cells = append(cells, cell)

				// Set cell as not frozen for the moment
				cell.Frozen = false
			}
		}

		// Shuffle cell traversal order
		rand.Shuffle(len(cells), func(i, j int) { cells[i], cells[j] = cells[j], cells[i] })

		// remove solution progressively
		var extractedCell Cell
		targetCellCount := grid.GetSize()*grid.GetSize() - emptyCellCount
		counter := 0
		for counter <= 1 && len(cells) > targetCellCount {
			extractedCell, cells = cells[0], cells[1:]
			grid.GetCell(extractedCell.Row, extractedCell.Column).Value = empty
			counter = grid.CountSolution()

			if counter > 1 {
				grid.GetCell(extractedCell.Row, extractedCell.Column).Value = extractedCell.Value
			} else {
				counter = 0
			}
		}

		if len(cells) == targetCellCount {
			for rowID, columns := range grid.Cells {
				for columnID := range columns {
					cell := grid.GetCell(rowID, columnID)

					// Set non empty cell as frozen
					if !cell.isEmpty() {
						cell.Frozen = true
					}
				}
			}

			return true
		}
	}

	return false
}
