package users

import (
	"sync"

	"github.com/google/uuid"
)

// User structure, with UUID and name
type User struct {
	UUID uuid.UUID
	Name string
}

// Users mapped to their UUID
type Users map[uuid.UUID]User

// LockableUsers are Users with a safe RWMutex
type LockableUsers struct {
	data Users
	sync.RWMutex
}
