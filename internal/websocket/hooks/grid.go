package hook

import (
	"log"

	"github.com/dezounet/sugokud/internal/websocket"
	"github.com/dezounet/sugokud/pkg/sudoku"
	"github.com/google/uuid"
)

// UpdateMessage sent to sync client's grid
type UpdateMessage struct {
	sudoku.Cell
	UUID uuid.UUID
}

// CreateGridHook instanciate and return hook for a sudoku.Grid,
// binding it to the input argument
func CreateGridHook(grid *sudoku.Grid) websocket.Hook {
	return websocket.Hook{
		OnConnection: nil,
		OnClose:      nil,
		Events: websocket.Events{
			websocket.GRIDUPDATE: struct{}{},
		},
		OnEvent: func(msg *websocket.Message) *websocket.Message {
			return onGridEvent(grid, msg)
		},
	}
}

func onGridEvent(grid *sudoku.Grid, msg *websocket.Message) *websocket.Message {
	var returnMsg *websocket.Message
	cell, ok := msg.Data.(UpdateMessage)
	if ok {
		switch msg.Event {
		case websocket.GRIDUPDATE:
			if cell.Row >= 0 && cell.Row < grid.GetSize() &&
				cell.Column >= 0 && cell.Column < grid.GetSize() &&
				cell.Value >= 0 && cell.Value <= grid.GetSize() &&
				!cell.Frozen && !grid.GetCell(cell.Row, cell.Column).Frozen &&
				grid.UUID == cell.UUID {
				// Update only if it is permitted
				returnMsg = msg
			} else {
				log.Printf("error[%s]: message is invalid - %v", msg.Event, cell)
			}
		default:
			log.Printf("error[%s]: event not handled", msg.Event)
		}
	} else {
		log.Printf("error[%s]: malformed message - %v", msg.Event, msg.Data)
	}

	return returnMsg
}
