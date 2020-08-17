package sudoku

import "testing"

func TestUsedInRow(t *testing.T) {
	grid := CreateEmptyGrid(2)

	for i := 1; i <= grid.GetSize(); i++ {
		if grid.usedInRow(0, i) || grid.usedInRow(1, i) || grid.usedInRow(2, i) || grid.usedInRow(3, i) {
			t.Fatal(i, "should not be detected as used in row, but it is")
		}
	}

	grid.Cells[0][0].Value = 1
	grid.Cells[1][1].Value = 1
	grid.Cells[2][2].Value = 1
	grid.Cells[3][3].Value = 1

	for i := 0; i < grid.GetSize(); i++ {
		if !grid.usedInRow(i, 1) || !grid.usedInRow(i, 1) || !grid.usedInRow(i, 1) || !grid.usedInRow(i, 1) {
			t.Fatal(i, "should be detected as used in row, and is not")
		}
	}
}

func TestUsedInColumn(t *testing.T) {
	grid := CreateEmptyGrid(2)

	for i := 1; i <= grid.GetSize(); i++ {
		if grid.usedInColumn(0, i) || grid.usedInColumn(1, i) || grid.usedInColumn(2, i) || grid.usedInColumn(3, i) {
			t.Fatal(i, "should not be detected as used in column, but it is")
		}
	}

	grid.Cells[0][0].Value = 1
	grid.Cells[1][1].Value = 1
	grid.Cells[2][2].Value = 1
	grid.Cells[3][3].Value = 1

	for i := 0; i < grid.GetSize(); i++ {
		if !grid.usedInColumn(i, 1) || !grid.usedInColumn(i, 1) || !grid.usedInColumn(i, 1) || !grid.usedInColumn(i, 1) {
			t.Fatal(i, "should be detected as used in column, and is not")
		}
	}
}

func TestUsedInBlock(t *testing.T) {
	grid := CreateEmptyGrid(2)

	for i := 1; i <= grid.GetSize(); i++ {
		if grid.usedInBlock(0, 0, i) || grid.usedInBlock(0, 1, i) || grid.usedInBlock(1, 0, i) || grid.usedInBlock(1, 1, i) {
			t.Fatal(i, "should not be detected as used in block, but it is")
		}
	}

	grid.Cells[0][0].Value = 1
	grid.Cells[0][1].Value = 2
	grid.Cells[1][0].Value = 3
	grid.Cells[1][1].Value = 4

	for i := 1; i <= grid.GetSize(); i++ {
		if !grid.usedInBlock(0, 0, i) || !grid.usedInBlock(0, 1, i) || !grid.usedInBlock(1, 0, i) || !grid.usedInBlock(1, 1, i) {
			t.Fatal(i, "should be detected as used in block, and is not")
		}
	}
}

func TestIsSafe(t *testing.T) {
	grid := CreateEmptyGrid(2)

	if grid.IsSafe(0, 0, 0) {
		t.Fatal(empty, "should not be detected as safe to use, but it is not")
	}

	if !grid.IsSafe(0, 0, 1) {
		t.Fatal("1 should be detected as safe to use, but it is not")
	}

	// row
	grid.Cells[0][2].Value = 1
	if grid.IsSafe(0, 0, 1) {
		t.Fatal("1 should not be detected as safe to use, but it is not")
	}

	grid.Cells[0][2].Value = 0

	// column
	grid.Cells[2][0].Value = 1
	if grid.IsSafe(0, 0, 1) {
		t.Fatal("1 should not be detected as safe to use, but it is not")
	}

	grid.Cells[0][2].Value = 1

	// block
	grid.Cells[1][1].Value = 1
	if grid.IsSafe(0, 0, 1) {
		t.Fatal("1 should not be detected as safe to use, but it is not")
	}
}
