package sudoku

import "testing"

func TestIsSolved(t *testing.T) {
	grid := CreateEmptyGrid(2)

	if grid.IsSolved() {
		t.Fatal("Grid should not be solved, but it is")
	}

	grid.Cells[0][0].Value = 1
	grid.Cells[0][1].Value = 2
	grid.Cells[0][2].Value = 3
	grid.Cells[0][3].Value = 4

	if grid.IsSolved() {
		t.Fatal("Grid should not be solved, but it is")
	}

	grid.Cells[1][0].Value = 3
	grid.Cells[1][1].Value = 4
	grid.Cells[1][2].Value = 1
	grid.Cells[1][3].Value = 2

	if grid.IsSolved() {
		t.Fatal("Grid should not be solved, but it is")
	}

	grid.Cells[2][0].Value = 2
	grid.Cells[2][1].Value = 3
	grid.Cells[2][2].Value = 4
	grid.Cells[2][3].Value = 1

	if grid.IsSolved() {
		t.Fatal("Grid should not be solved, but it is")
	}

	grid.Cells[3][0].Value = 4
	grid.Cells[3][1].Value = 1
	grid.Cells[3][2].Value = 2

	if grid.IsSolved() {
		t.Fatal("Grid should not be solved, but it is")
	}

	grid.Cells[3][3].Value = 3

	if !grid.IsSolved() {
		t.Fatal("Grid should be solved, but it is not")
	}
}

func TestSolve(t *testing.T) {
	grid := CreateEmptyGrid(1)

	if !grid.Solve() {
		t.Fatal("Grid should be solved, but it is not")
	}

	if !grid.IsSolved() {
		t.Fatal("Grid should be solved, but it is not")
	}
}

func TestCountSolution(t *testing.T) {
	grid := CreateEmptyGrid(1)
	counter := grid.CountSolution()
	if counter != 1 {
		t.Fatal("1 solution expected, got", counter)
	}

	grid = CreateEmptyGrid(2)
	grid.Cells[0][0].Value = 1
	grid.Cells[0][1].Value = 2
	grid.Cells[0][2].Value = 3
	grid.Cells[0][3].Value = 4
	grid.Cells[1][0].Value = 3
	grid.Cells[1][1].Value = 4
	grid.Cells[1][2].Value = 1
	grid.Cells[1][3].Value = 2
	counter = grid.CountSolution()
	if counter != 4 {
		t.Fatal("4 solutions expected, got", counter)
	}

	grid.Cells[2][0].Value = 2
	grid.Cells[2][1].Value = 3
	grid.Cells[2][2].Value = 4
	grid.Cells[2][3].Value = 1
	counter = grid.CountSolution()
	if counter != 1 {
		t.Fatal("1 solution expected, got", counter)
	}

	grid.Cells[3][0].Value = 4
	grid.Cells[3][1].Value = 1
	grid.Cells[3][2].Value = 2
	grid.Cells[3][3].Value = 3
	counter = grid.CountSolution()
	if counter != 1 {
		t.Fatal("1 solution expected, got", counter)
	}

}
