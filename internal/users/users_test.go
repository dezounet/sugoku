package users

import (
	"strconv"
	"testing"

	"github.com/google/uuid"
)

func TestCreate(t *testing.T) {
	users := Create()

	if users.Count() != 0 {
		t.Error("Users created should be empty")
	}
}

func TestCount(t *testing.T) {
	users := Create()
	uuids := make([]string, 10)
	name := "test"

	for i := 0; i < 10; i++ {
		uuids[i] = uuid.New().String()
	}

	userCount := users.Count()
	if userCount != 0 {
		t.Error("Users should contain 0 entry, got", userCount)
	}

	for i := 1; i < 10; i++ {
		user := User{
			UUID: uuids[i-1],
			Name: name + strconv.Itoa(i),
		}
		users.AddOrUpdate(user)

		userCount := users.Count()
		if userCount != i {
			t.Error("Users should contain", i, "entries, got", userCount)
		}
	}
}

func TestGet(t *testing.T) {
	users := Create()

	uuids := make([]string, 10)
	names := make([]string, 10)
	for i := 0; i < 10; i++ {
		uuids[i] = uuid.New().String()
		names[i] = "test" + strconv.Itoa(i)
		user := User{
			UUID: uuids[i],
			Name: names[i],
		}
		users.AddOrUpdate(user)
	}

	for i := 0; i < 10; i++ {
		user, ok := users.Get(uuids[i])

		if !ok {
			t.Error("User", uuids[i], "should have been found")
		}

		if user.UUID != uuids[i] || user.Name != names[i] {
			t.Error("User", uuids[i], "not retrieved as expected, got", user)
		}
	}

	for i := 0; i < 10; i++ {
		unknownUUID := uuid.New().String()
		if _, ok := users.Get(unknownUUID); ok {
			t.Error("User", unknownUUID, "shouldn't have been found")
		}
	}
}

func TestAdd(t *testing.T) {
	users := Create()
	uuid := uuid.New().String()
	name := "test"
	user := User{
		UUID: uuid,
		Name: name,
	}

	users.AddOrUpdate(user)

	userCount := users.Count()
	if userCount != 1 {
		t.Error("Users should contain 1 entry, got", userCount)
	}

	retrievedUser, ok := users.Get(uuid)
	if !ok {
		t.Error("User should be in list")
	}
	if retrievedUser.Name != name {
		t.Error("User name should be", name, "got", retrievedUser.Name)
	}
}

func TestUpdate(t *testing.T) {
	users := Create()
	uuid := uuid.New().String()
	name := "test "

	user := User{
		UUID: uuid,
		Name: name,
	}

	users.AddOrUpdate(user)

	newName := "new"
	user.Name = newName

	users.AddOrUpdate(user)

	userCount := users.Count()
	if userCount != 1 {
		t.Error("Users should contain 1 entry, got", userCount)
	}

	user, ok := users.Get(uuid)
	if !ok {
		t.Error("User should be in list")
	}
	if user.UUID != uuid || user.Name != newName {
		t.Error("User", uuid, "not retrieved as expected, got", user)
	}
}

func TestRemove(t *testing.T) {
	users := Create()
	knownUUID := uuid.New().String()
	name := "test "

	user := User{
		UUID: knownUUID,
		Name: name,
	}

	users.AddOrUpdate(user)

	unknownUUUID := uuid.New().String()
	users.Remove(unknownUUUID)

	userCount := users.Count()
	if userCount != 1 {
		t.Error("Users should contain 1 entry, got", userCount)
	}

	_, ok := users.Get(unknownUUUID)
	if ok {
		t.Error("Unknown user shouldn't be in list")
	}

	users.Remove(knownUUID)

	userCount = users.Count()
	if userCount != 0 {
		t.Error("Users should contain 0 entry, got", userCount)
	}

	user, ok = users.Get(knownUUID)
	if ok {
		t.Error("Known user should be in list")
	}

}

func TestSample(t *testing.T) {
	users := Create()

	count := 100
	uuids := make([]string, count)

	for i := 0; i < count; i++ {
		uuids[i] = uuid.New().String()
	}

	for i := 1; i < count; i++ {
		user := User{
			UUID: uuids[i-1],
			Name: strconv.Itoa(i),
		}
		users.AddOrUpdate(user)
	}

	sampleCount := uint8(10)
	samples := users.Sample(sampleCount)
	if len(samples) != int(sampleCount) {
		t.Error("Expected a sample a size", sampleCount, "got", len(samples))
	}
}

func TestEmptySample(t *testing.T) {
	users := Create()
	sampleCount := uint8(10)
	samples := users.Sample(sampleCount)
	if len(samples) != 0 {
		t.Error("Expected a sample a size of 0, got", len(samples))
	}
}
