package sudoku

import "testing"

func TestDifficultyToEmptyCellCount(t *testing.T) {
	grid := CreateEmptyGrid(3)

	easyCount := grid.DifficultyToEmptyCellCount(EASY)
	mediumCount := grid.DifficultyToEmptyCellCount(MEDIUM)
	hardCount := grid.DifficultyToEmptyCellCount(HARD)
	nightmareCount := grid.DifficultyToEmptyCellCount(NIGHTMARE)

	if easyCount >= mediumCount || mediumCount >= hardCount || hardCount >= nightmareCount {
		t.Fatal("Expecting difficulty to increase, but it is not")
	}

	cellCount := grid.GetSize() * grid.GetSize()
	if easyCount > cellCount || mediumCount > cellCount || hardCount > cellCount || nightmareCount > cellCount {
		t.Fatal("There can't be more empty cell than total cell count")
	}
}

func TestInitialize(t *testing.T) {
	grid := CreateEmptyGrid(2)
	cellCount := grid.GetSize() * grid.GetSize()

	grid.Initialize(NIGHTMARE)

	if grid.CountEmptyCell() == cellCount || grid.CountEmptyCell() == 0 || !grid.Solve() {
		t.Fatal("Initialized grid should contain empty cell!")
	}
}

func TestRemoveCellValue(t *testing.T) {
	grid := CreateEmptyGrid(1)
	grid.GetCell(0, 0).Value = 1

	if !grid.RemoveCellValue(1) || !grid.GetCell(0, 0).isEmpty() {
		t.Fatal("Grid should have been emptied, but it is not")
	}

	grid = CreateEmptyGrid(2)
	if !grid.RemoveCellValue(9) {
		t.Fatal("Grid should have been emptied, but it is not")
	}
}

func TestFrozenRemoveCellValue(t *testing.T) {
	grid := CreateEmptyGrid(2)
	expectedEmptyCell := 8
	if !grid.RemoveCellValue(expectedEmptyCell) {
		t.Fatal("Grid should have been emptied, but it is not")
	}

	frozenCellCount := 0
	for rowID, columns := range grid.Cells {
		for columnID := range columns {
			cell := grid.GetCell(rowID, columnID)
			if cell.isEmpty() {
				if cell.Frozen {
					t.Fatal("Empty cell should not be frozen")
				}
			}

			if cell.Frozen {
				frozenCellCount++
			}
		}
	}

	if frozenCellCount != expectedEmptyCell {
		t.Fatal("Expecting", expectedEmptyCell, " frozen cells, but go", frozenCellCount)
	}
}

func TestRemoveTooMuchCellValue(t *testing.T) {
	grid := CreateEmptyGrid(1)
	grid.GetCell(0, 0).Value = 1

	if grid.RemoveCellValue(2) || grid.GetCell(0, 0).Value != 1 {
		t.Fatal("Grid should not have been emptied, but it is")
	}
}

func TestUnsolvableRemoveCellValue(t *testing.T) {
	grid := CreateEmptyGrid(2)
	grid.GetCell(0, 0).Value = 1
	grid.GetCell(0, 1).Value = 1

	if grid.RemoveCellValue(1) {
		t.Fatal("Grid is unsolvable and should not have been emptied, but it is")
	}
}
