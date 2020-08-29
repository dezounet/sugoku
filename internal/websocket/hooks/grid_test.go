package hook

import (
	"testing"

	"github.com/dezounet/sugokud/internal/websocket"
	"github.com/dezounet/sugokud/pkg/sudoku"
)

func TestGridEvents(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	hook := CreateGridHook(grid)

	expectedEvents := []websocket.Event{
		websocket.GRIDUPDATE,
		websocket.GRIDRESET,
	}

	for _, event := range expectedEvents {
		if _, ok := hook.Events[event]; !ok {
			t.Fatal("Expected event", event, "but it was not found")
		}
	}
}

func TestOnGridUpdateEvent(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	hook := CreateGridHook(grid)

	testMsg := websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    0,
					Column: 0,
				},
				Value:  0,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg := hook.OnEvent(&testMsg)
	if msg == nil {
		t.Fatal("Expecting non nil message")
	}

	if *msg != testMsg {
		t.Fatal("Expecting message forwarded to be the same, but got", msg)
	}
}

func TestOnGridUpdateEventInvalidRow(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	hook := CreateGridHook(grid)

	testMsg := websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    9,
					Column: 0,
				},
				Value:  1,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg := hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}

	testMsg = websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    -1,
					Column: 0,
				},
				Value:  1,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg = hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}
}

func TestOnGridUpdateEventInvalidColumn(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	hook := CreateGridHook(grid)

	testMsg := websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    0,
					Column: 9,
				},
				Value:  1,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg := hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}

	testMsg = websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    0,
					Column: -1,
				},
				Value:  1,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg = hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}
}

func TestOnGridUpdateEventInvalidValue(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	hook := CreateGridHook(grid)

	testMsg := websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    0,
					Column: 0,
				},
				Value:  -1,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg := hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}

	testMsg = websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    0,
					Column: 0,
				},
				Value:  10,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg = hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}
}

func TestOnGridUpdateEventOnFrozenCell(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	grid.GetCell(0, 0).Frozen = true

	hook := CreateGridHook(grid)

	testMsg := websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    0,
					Column: 0,
				},
				Value:  1,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg := hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}

	testMsg = websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    1,
					Column: 0,
				},
				Value:  1,
				Frozen: true,
			},
			UUID: grid.UUID,
		},
	}

	msg = hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}
}

func TestGridUnknownEvent(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	hook := CreateGridHook(grid)

	testMsg := websocket.Message{
		Event: "unknown",
		Data: UpdateMessage{
			Cell: sudoku.Cell{
				Coordinates: sudoku.Coordinates{
					Row:    0,
					Column: 0,
				},
				Value:  0,
				Frozen: false,
			},
			UUID: grid.UUID,
		},
	}

	msg := hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message on second connection")
	}
}

func TestGridMalformedMessage(t *testing.T) {
	grid := sudoku.CreateEmptyGrid(3)
	hook := CreateGridHook(grid)

	testMsg := websocket.Message{
		Event: websocket.GRIDUPDATE,
		Data:  "malformed data",
	}

	msg := hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message on second connection")
	}
}
