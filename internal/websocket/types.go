package websocket

import (
	"sync"

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
	Event Event       `json:"event"`
	UUID  string      `json:"uuid,omitempty"`
	Data  interface{} `json:"data"`
}

// LockableConnectionHandler container with a safe RWMutex
// A connection is indexed by its websocket, and ssociated
// with a UUID
type LockableConnectionHandler struct {
	sync.RWMutex
	Upgrader  websocket.Upgrader
	Broadcast chan Message
	M         map[*websocket.Conn]string
	Hooks     []Hook
}

// Hook to be called by LockableConnectionHandler in websocket lifecycle
type Hook struct {
	OnConnection func(string) *Message
	OnClose      func(string) *Message
	Events       Events
	OnEvent      func(*Message) *Message
}
