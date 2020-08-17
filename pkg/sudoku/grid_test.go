package sudoku

import (
	"testing"
)

func TestIsEmpty(t *testing.T) {
	cell := Cell{}

	if !cell.isEmpty() {
		t.Fatal("Cell should be flagged as empty, and is not")
	}

	for i := 1; i < 100; i++ {
		cell.Value = i
		if cell.isEmpty() {
			t.Fatal("Cell should be flagged as not empty, and is")
		}
	}
}

func TestGetSize(t *testing.T) {
	grid := CreateEmptyGrid(3)

	if grid.GetSize() != 9 {
		t.Fatal("Grid size should be 9, got", grid.GetSize())
	}

	grid = CreateEmptyGrid(2)

	if grid.GetSize() != 4 {
		t.Fatal("Grid size should be 2, got", grid.GetSize())
	}
}

func TestGetCell(t *testing.T) {
	grid := CreateEmptyGrid(3)

	for i := 0; i < grid.GetSize(); i++ {
		for j := 0; j < grid.GetSize(); j++ {
			cell := &grid.Cells[i][j]

			if grid.GetCell(i, j) != cell {
				t.Fatal("Cell", i, "x", j, "should match, but did not")
			}
		}
	}
}

func TestCreateEmptyGrid(t *testing.T) {
	grid := CreateEmptyGrid(3)

	for rowID, columns := range grid.Cells {
		for columnID, cell := range columns {
			if cell.Value != 0 || cell.Row != rowID || cell.Column != columnID || cell.Frozen {
				t.Fatal("Wrong empty initialization for a grid")
			}
		}
	}
}

func TestReset(t *testing.T) {
	grid := CreateEmptyGrid(3)

	grid.Cells[0][0].Value = 1
	grid.Cells[0][0].Frozen = true

	grid.Reset()

	if grid.Cells[0][0].Value != 0 || grid.Cells[0][0].Frozen {
		t.Fatal("Expecting grid to be reset, but it is not:", grid.Cells[0][0])
	}
}

func TestFindEmptyCell(t *testing.T) {
	grid := CreateEmptyGrid(3)

	row := -1
	column := -1

	if !grid.findEmptyCell(&row, &column) {
		t.Fatal("Expecting empty cell, got none")
	}
	if row != 0 || column != 0 {
		t.Fatal("Expected cell 0x0, got", row, "x", column)
	}

	grid.Cells[0][0].Value = 1
	if !grid.findEmptyCell(&row, &column) {
		t.Fatal("Expecting empty cell, got none")
	}
	if row != 0 || column != 1 {
		t.Fatal("Expected cell 0x1, got", row, "x", column)
	}
}

func TestFindNoEmptyCell(t *testing.T) {
	grid := CreateEmptyGrid(1)
	grid.Cells[0][0].Value = 1

	row := -1
	column := -1

	if grid.findEmptyCell(&row, &column) {
		t.Fatal("Expecting no empty cell, got one")
	}
}

func TestCountEmptyCell(t *testing.T) {
	for i := 1; i < 4; i++ {
		grid := CreateEmptyGrid(i)
		cellCount := grid.GetSize() * grid.GetSize()

		if grid.CountEmptyCell() != cellCount {
			t.Fatal("Expecting", cellCount, "empty cell, got", grid.CountEmptyCell())
		}

		grid.Cells[0][0].Value = 1

		if grid.CountEmptyCell() != cellCount-1 {
			t.Fatal("Expecting", cellCount-1, "empty cell, got", grid.CountEmptyCell())
		}
	}
}
