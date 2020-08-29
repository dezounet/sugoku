package websocket

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// messageSizeLimit of received messages. Connection is closed upon
// bigger message
var messageSizeLimit = int64(1024)

// getUpgrader instanciate and return an upgrader object, that should be
// used to upgrade a basic connection to a websocket connection
// No specific option is passed to this upgrader, modify here if you want
// to adapt the Upgrader behaviour
func getUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		// Uncomment these lines to allow request from any location
		// CheckOrigin: func(r *http.Request) bool {
		// return true
		// },
	}
}

// GetHandler instanciate and return an http request handler to promote
// the connection to a websocket.
func GetHandler() LockableConnectionHandler {
	return LockableConnectionHandler{
		Upgrader:  getUpgrader(),
		Broadcast: make(chan Message),
		M:         make(map[*websocket.Conn]string),
		Hooks:     make([]Hook, 0),
	}
}

func (h *LockableConnectionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get optionnal parameter uuid, or generate one
	uuidParameter := r.URL.Query().Get("uuid")
	UUID, err := uuid.Parse(uuidParameter)
	if err != nil {
		UUID = uuid.New()
	}

	// Upgrade initial GET request to a websocket
	ws, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	} else {
		h.Lock()
		h.M[ws] = UUID.String()
		h.Unlock()

		// Limit max message size to prevent DoS
		ws.SetReadLimit(messageSizeLimit)
	}

	// Make sure we close the connection when the function returns
	defer closeConnection(h, ws)

	// Iterate over hooks to perform OnConnection actions
	for _, hook := range h.Hooks {
		if hook.OnConnection != nil {
			msg := hook.OnConnection(UUID.String())
			if msg != nil {

				h.Broadcast <- *msg
			}
		}
	}

	for {
		// Read in a new message as JSON and map it to a Message object
		// var msg Message
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			// Close connection on error
			log.Printf("error[%s]: %v", UUID, err)
			return
		}

		// Set current user
		msg.UUID = UUID.String()
		h.Broadcast <- msg
	}
}

func closeConnection(h *LockableConnectionHandler, ws *websocket.Conn) {
	if ws != nil {
		// Close connection to prevent any new message from this client
		ws.Close()

		// Get client UUID
		h.RLock()
		UUID := h.M[ws]
		h.RUnlock()

		// Delete the corresponding websocket entry
		h.Lock()
		delete(h.M, ws)
		h.Unlock()

		// Iterate over hooks to perform OnClose actions
		for _, hook := range h.Hooks {
			if hook.OnConnection != nil {
				msg := hook.OnClose(UUID)
				if msg != nil {
					h.Broadcast <- *msg
				}
			}
		}
	}
}

// HandleMessages received from the handler boradcast channel, forward
// it to connected clients
func HandleMessages(handler *LockableConnectionHandler) {
	for {
		// Grab the next message from the broadcast channel
		msg := <-handler.Broadcast

		// Trigger event processing for every hooks
		log.Println("broadcasting", msg)
		for _, hook := range handler.Hooks {
			// Only trigger if the current hook is interested in this event
			if _, ok := hook.Events[msg.Event]; ok {
				hookMsg := hook.OnEvent(&msg)
				if hookMsg != nil {
					// Send it out to every client that is currently connected
					broadcastMessage(handler, hookMsg)
				}
			}
		}
	}
}

// RegisterHook to a LockableConnectionHandler instance
func RegisterHook(handler *LockableConnectionHandler, hook Hook) {
	handler.Lock()
	handler.Hooks = append(handler.Hooks, hook)
	handler.Unlock()
}

// broadcastMessage to every connected client
// Disconnet client on error
func broadcastMessage(handler *LockableConnectionHandler, msg *Message) {
	handler.RLock()
	for ws, UUID := range handler.M {
		err := ws.WriteJSON(*msg)
		if err != nil {
			log.Printf("error[%s]: %v", UUID, err)
			closeConnection(handler, ws)
		}
	}
	handler.RUnlock()
}
