package hook

import (
	"log"

	"github.com/dezounet/sugokud/internal/websocket"
	"github.com/dezounet/sugokud/pkg/sudoku"
	"github.com/mitchellh/mapstructure"
)

// UpdateMessage sent to sync client's grid
type UpdateMessage struct {
	sudoku.Cell `mapstructure:",squash"`
	UUID        string `json:"uuid"`
}

// CreateGridHook instanciate and return hook for a sudoku.Grid,
// binding it to the input argument
func CreateGridHook(grid *sudoku.Grid) websocket.Hook {
	return websocket.Hook{
		OnConnection: nil,
		OnClose:      nil,
		Events: websocket.Events{
			websocket.GRIDUPDATE: struct{}{},
			websocket.GRIDRESET:  struct{}{},
		},
		OnEvent: func(msg *websocket.Message) *websocket.Message {
			return onGridEvent(grid, msg)
		},
	}
}

func onGridEvent(grid *sudoku.Grid, msg *websocket.Message) *websocket.Message {
	var returnMsg *websocket.Message

	var cell UpdateMessage
	err := mapstructure.Decode(msg.Data, &cell)
	if err == nil {
		switch msg.Event {
		case websocket.GRIDUPDATE:
			// Update only if it is permitted
			if cell.Row >= 0 && cell.Row < grid.GetSize() &&
				cell.Column >= 0 && cell.Column < grid.GetSize() &&
				cell.Value >= 0 && cell.Value <= grid.GetSize() &&
				!cell.Frozen && !grid.GetCell(cell.Row, cell.Column).Frozen &&
				grid.UUID == cell.UUID {
				// save state locally
				grid.GetCell(cell.Row, cell.Column).Value = cell.Value

				// broadcast message
				returnMsg = msg
			} else {
				log.Printf("error[%s]: message is invalid - %v | %s", msg.Event, cell, grid.UUID)
			}
		case websocket.GRIDRESET:
			returnMsg = msg
		default:
			log.Printf("error[%s]: event not handled", msg.Event)
		}
	} else {
		log.Printf("error[%s]: malformed message - %v because %s", msg.Event, msg.Data, err)
	}

	return returnMsg
}
