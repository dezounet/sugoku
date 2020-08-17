package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/dezounet/sugokud/internal/websocket"
	"github.com/dezounet/sugokud/pkg/sudoku"
	"github.com/google/uuid"
)

func TestGetEmptyGrid(t *testing.T) {
	for size := 0; size <= 5; size++ {
		grid := sudoku.CreateEmptyGrid(size)
		handler := GetGridHandler{
			Grid: grid,
		}

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, nil)

		var data sudoku.Grid
		err := json.Unmarshal(rec.Body.Bytes(), &data)
		if err != nil {
			t.Fatal("Failed to understand response:", rec.Body.String())
		}

		// Check we got back the right block size
		if data.BlockSize != size {
			t.Fatal("Expected Block size of ", size, "but got", grid.BlockSize)
		}

		// Check we got back the right UUID
		if data.UUID != grid.UUID {
			t.Fatal("Expected UUID ", data.UUID, "but got", grid.UUID)
		}

		// Check we got back the right cell count
		if data.CountEmptyCell() != (grid.GetSize() * grid.GetSize()) {
			t.Fatal("Expected", grid.GetSize()*grid.GetSize(), "empty cells, got", data.CountEmptyCell())
		}

		// Check we got back the right row / column count
		rowCount := 0
		for _, columns := range data.Cells {
			rowCount++

			columnCount := 0
			for range columns {
				columnCount++
			}

			if columnCount != grid.GetSize() {
				t.Fatal("Expected", grid.GetSize(), "columns, got", columnCount)
			}
		}
		if rowCount != grid.GetSize() {
			t.Fatal("Expected", grid.GetSize(), "rows, got", rowCount)
		}

	}
}

func TestGetGridUUID(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	handler := GetGridUUIDHandler{
		Grid: grid,
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, nil)

	var data uuid.UUID
	err := json.Unmarshal(rec.Body.Bytes(), &data)
	if err != nil {
		t.Fatal("Failed to understand response:", rec.Body.String())
	}

	if data != grid.UUID {
		t.Fatal("Expected UUID ", data, "but got", grid.UUID)
	}
}

func TestResetGrid(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	grid.Solve()

	UUID := grid.UUID

	handler := GetGridResetHandler{
		Grid: grid,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com?difficulty=easy", nil)
	handler.ServeHTTP(rec, req)

	if grid.UUID == UUID {
		t.Fatal("Expected new UUID, but no change detected")
	}
}

func TestResetNonSolvedGrid(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	UUID := grid.UUID
	handler := GetGridResetHandler{
		Grid: grid,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com?difficulty=easy", nil)
	handler.ServeHTTP(rec, req)

	if grid.UUID != UUID {
		t.Fatal("Expected same UUID, change detected")
	}
}

func TestBroadcastReset(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(1)
	grid.Solve()

	handler := GetGridResetHandler{
		Grid:      grid,
		Broadcast: make(chan websocket.Message),
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com?difficulty=easy", nil)
	go handler.ServeHTTP(rec, req)

	msg := <-handler.Broadcast

	if msg.Event != websocket.GRIDRESET {
		t.Fatal("Expecting GRIDRESET event, got", msg.Event)
	}
}

func TestIncreaseDifficulty(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	grid.Solve()

	handler := GetGridResetHandler{
		Grid: grid,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com?difficulty=easy", nil)
	handler.ServeHTTP(rec, req)
	easyCount := grid.CountEmptyCell()

	grid.Solve()

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://test.com?difficulty=medium", nil)
	handler.ServeHTTP(rec, req)
	mediumCount := grid.CountEmptyCell()

	grid.Solve()

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://test.com?difficulty=hard", nil)
	handler.ServeHTTP(rec, req)
	hardCount := grid.CountEmptyCell()

	grid.Solve()

	rec = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "http://test.com?difficulty=nightmare", nil)
	handler.ServeHTTP(rec, req)
	nightmareCount := grid.CountEmptyCell()

	if easyCount >= mediumCount || mediumCount >= hardCount || hardCount >= nightmareCount {
		t.Fatal("Difficulty level not enforced as expected")
	}
}

func TestUnknownDifficulty(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	grid.Solve()
	handler := GetGridResetHandler{
		Grid: grid,
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test.com?difficulty=qwerty", nil)
	handler.ServeHTTP(rec, req)
	emptyCellCount := grid.CountEmptyCell()

	if emptyCellCount == 0 || emptyCellCount == grid.GetSize()*grid.GetSize() {
		t.Fatal("Difficulty level not enforced as expected, got ", emptyCellCount, "empty cell(s)")
	}
}
