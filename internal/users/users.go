package users

import (
	"github.com/google/uuid"
)

// Create an empty instance of LockableUsers
func Create() LockableUsers {
	return LockableUsers{
		data: make(Users),
	}
}

// AddOrUpdate a User entry inside the users
func (users *LockableUsers) AddOrUpdate(user User) {
	users.Lock()
	users.data[user.UUID] = user
	users.Unlock()
}

// Remove a User from the users
func (users *LockableUsers) Remove(UUID uuid.UUID) {
	users.Lock()
	delete(users.data, UUID)
	users.Unlock()
}

// Get a specific User, or nil if it does not exist
func (users *LockableUsers) Get(UUID uuid.UUID) (User, bool) {
	users.RLock()
	user, ok := users.data[UUID]
	users.RUnlock()

	return user, ok
}

// Sample of the users
func (users *LockableUsers) Sample(count uint8) []User {
	users.RLock()
	userCount := len(users.data)
	users.RUnlock()

	if int(count) > userCount {
		count = uint8(userCount)
	}

	samples := make([]User, count)

	users.RLock()
	i := uint8(0)
	for _, user := range users.data {
		if i >= count {
			break
		}

		samples[i] = user
		i++
	}
	users.RUnlock()

	return samples
}

// Count the users
func (users *LockableUsers) Count() int {
	users.Lock()
	userCount := len(users.data)
	users.Unlock()

	return userCount
}
