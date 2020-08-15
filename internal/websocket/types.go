package websocket

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Event to discriminate message received
type Event string

// Events as a set
type Events map[Event]struct{}

// Different message type available
const (
	USERADD    Event = "user-add"
	USERDEL          = "user-del"
	USERUPDATE       = "user-update"
	GRIDRESET        = "grid-reset"
	GRIDUPDATE       = "grid-update"
)

// Message structure to share information
type Message struct {
	Event Event
	Data  interface{}
}

// LockableConnectionHandler container with a safe RWMutex
// A connection is indexed by its websocket, and ssociated
// with a UUID
type LockableConnectionHandler struct {
	sync.RWMutex
	Upgrader  websocket.Upgrader
	Broadcast chan Message
	M         map[*websocket.Conn]uuid.UUID
	Hooks     []Hook
}

// Hook to be called by LockableConnectionHandler in websocket lifecycle
type Hook struct {
	OnConnection func(uuid.UUID) *Message
	OnClose      func(uuid.UUID) *Message
	Events       Events
	OnEvent      func(*Message) *Message
}
