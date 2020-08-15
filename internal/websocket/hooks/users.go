package hook

import (
	"log"

	"github.com/dezounet/sugokud/internal/users"
	"github.com/dezounet/sugokud/internal/websocket"
	"github.com/google/uuid"
)

// CreateUserHook instanciate and return hook for LockableConnectionHandler,
// binding it to the input argument
func CreateUserHook(users *users.LockableUsers) websocket.Hook {
	return websocket.Hook{
		OnConnection: func(UUID uuid.UUID) *websocket.Message {
			return addUser(users, UUID, "unnamed user")
		},
		OnClose: func(UUID uuid.UUID) *websocket.Message {
			return removeUser(users, UUID)
		},
		Events: websocket.Events{
			websocket.USERADD:    struct{}{},
			websocket.USERUPDATE: struct{}{},
			websocket.USERDEL:    struct{}{},
		},
		OnEvent: func(msg *websocket.Message) *websocket.Message {
			return onUserEvent(users, msg)
		},
	}
}

func addUser(lockableUsers *users.LockableUsers, UUID uuid.UUID, name string) *websocket.Message {
	var msg *websocket.Message

	// Add new user to users. If UUID is already known, do nothing
	user, ok := lockableUsers.Get(UUID)
	if !ok {
		user = users.User{
			UUID: UUID,
			Name: name,
		}

		lockableUsers.AddOrUpdate(user)

		msg = &websocket.Message{
			Event: websocket.USERADD,
			Data:  user,
		}
	}

	return msg
}

func removeUser(lockableUsers *users.LockableUsers, UUID uuid.UUID) *websocket.Message {
	var msg *websocket.Message

	_, ok := lockableUsers.Get(UUID)
	if ok {
		// Delete user from users
		lockableUsers.Remove(UUID)

		// Broadcast to other clients
		msg = &websocket.Message{
			Event: websocket.USERDEL,
			Data:  UUID,
		}
	}

	return msg
}

func onUserEvent(lockableUsers *users.LockableUsers, msg *websocket.Message) *websocket.Message {
	var returnMsg *websocket.Message
	user, ok := msg.Data.(users.User)
	if ok {
		switch msg.Event {
		case websocket.USERADD:
			returnMsg = addUser(lockableUsers, user.UUID, user.Name)
		case websocket.USERUPDATE:
			if _, ok := lockableUsers.Get(user.UUID); ok {
				lockableUsers.AddOrUpdate(user)
				returnMsg = msg
			}
		case websocket.USERDEL:
			returnMsg = removeUser(lockableUsers, user.UUID)
		default:
			log.Printf("error[%s]: event not handled", msg.Event)
		}
	} else {
		log.Printf("error[%s]: malformed message - %v", msg.Event, msg.Data)
	}

	return returnMsg
}
