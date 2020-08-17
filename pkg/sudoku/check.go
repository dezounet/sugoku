package sudoku

func (grid *Grid) usedInRow(row int, value int) bool {
	for i := 0; i < grid.GetSize(); i++ {
		if grid.GetCell(row, i).Value == value {
			return true
		}
	}

	return false
}

func (grid *Grid) usedInColumn(column int, value int) bool {
	for i := 0; i < grid.GetSize(); i++ {
		if grid.GetCell(i, column).Value == value {
			return true
		}
	}

	return false
}

func (grid *Grid) usedInBlock(row int, column int, value int) bool {
	boxRow := row / grid.BlockSize
	boxColumn := column / grid.BlockSize

	for i := 0; i < grid.BlockSize; i++ {
		for j := 0; j < grid.BlockSize; j++ {
			if grid.GetCell(grid.BlockSize*boxRow+i, grid.BlockSize*boxColumn+j).Value == value {
				return true
			}
		}
	}

	return false
}

// IsSafe checks if a value can safely be used in a cell,
// without violating any sudoku rule
func (grid *Grid) IsSafe(row int, column int, value int) bool {
	return value != empty &&
		!grid.usedInRow(row, value) &&
		!grid.usedInColumn(column, value) &&
		!grid.usedInBlock(row, column, value)
}
