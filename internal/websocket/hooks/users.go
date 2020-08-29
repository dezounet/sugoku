package hook

import (
	"log"

	"github.com/dezounet/sugokud/internal/users"
	"github.com/dezounet/sugokud/internal/websocket"
	"github.com/mitchellh/mapstructure"
)

// CreateUserHook instanciate and return hook for LockableConnectionHandler,
// binding it to the input argument
func CreateUserHook(users *users.LockableUsers) websocket.Hook {
	return websocket.Hook{
		OnConnection: func(UUID string) *websocket.Message {
			return addUser(users, UUID, "unnamed user")
		},
		OnClose: func(UUID string) *websocket.Message {
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

func addUser(lockableUsers *users.LockableUsers, UUID string, name string) *websocket.Message {
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

func removeUser(lockableUsers *users.LockableUsers, UUID string) *websocket.Message {
	var msg *websocket.Message

	user, ok := lockableUsers.Get(UUID)
	if ok {
		// Delete user from users
		lockableUsers.Remove(UUID)

		// Broadcast to other clients
		msg = &websocket.Message{
			Event: websocket.USERDEL,
			Data: users.User{
				UUID: user.UUID,
				Name: user.Name,
			},
		}
	}

	return msg
}

func onUserEvent(lockableUsers *users.LockableUsers, msg *websocket.Message) *websocket.Message {
	var returnMsg *websocket.Message

	var user users.User
	err := mapstructure.Decode(msg.Data, &user)
	if err == nil {
		switch msg.Event {
		case websocket.USERADD:
			returnMsg = addUser(lockableUsers, user.UUID, user.Name)
		case websocket.USERUPDATE:
			if _, ok := lockableUsers.Get(msg.UUID); ok {
				// make sure uuid is set by server
				user.UUID = msg.UUID
				msg.Data = user
				returnMsg = msg

				// save state locally
				lockableUsers.AddOrUpdate(user)
			} else {
				log.Println("Failed to update unknown user", msg.UUID)
			}
		case websocket.USERDEL:
			returnMsg = removeUser(lockableUsers, user.UUID)
		default:
			log.Printf("error[%s]: event not handled", msg.Event)
		}
	} else {
		log.Printf("error[%s]: malformed message - %v because %s", msg.Event, msg.Data, err)
	}

	return returnMsg
}
