package users

import (
	"sync"
)

// User structure, with UUID and name
type User struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// Users mapped to their UUID
type Users map[string]User

// LockableUsers are Users with a safe RWMutex
type LockableUsers struct {
	data Users
	sync.RWMutex
}
