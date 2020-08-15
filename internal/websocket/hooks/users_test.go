package hook

import (
	"testing"

	"github.com/dezounet/sugokud/internal/users"
	"github.com/dezounet/sugokud/internal/websocket"
	"github.com/google/uuid"
)

func TestEvents(t *testing.T) {
	users := users.Create()
	hook := CreateUserHook(&users)

	expectedEvents := []websocket.Event{
		websocket.USERADD,
		websocket.USERUPDATE,
		websocket.USERDEL,
	}

	for _, event := range expectedEvents {
		if _, ok := hook.Events[event]; !ok {
			t.Fatal("Expected event", event, "but it was not found")
		}
	}
}

func TestOnConnection(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()

	msg := hook.OnConnection(UUID)
	if msg == nil {
		t.Fatal("Expecting non nil message")
	}
	if msg.Event != websocket.USERADD {
		t.Fatal("Expecting event", websocket.USERADD, "got", msg.Event)
	}

	typedMessage, ok := msg.Data.(users.User)
	if !ok {
		t.Fatal("Expecting User structure, can't cast it:", msg.Data)
	}
	if typedMessage.UUID != UUID {
		t.Fatal("Expecting UUID", UUID, "got", typedMessage.UUID)
	}

	if count := testUsers.Count(); count != 1 {
		t.Fatal("Expecting 1 user, got", count)
	}

	user, ok := testUsers.Get(UUID)
	if !ok || user.UUID != UUID {
		t.Fatal("Failed to retrieve UUID", UUID, "in users")
	}
}

func TestOnSecondConnection(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()

	msg := hook.OnConnection(UUID)
	msg = hook.OnConnection(UUID)
	if msg != nil {
		t.Fatal("Expecting nil message on second connection")
	}

	if count := testUsers.Count(); count != 1 {
		t.Fatal("Expecting 1 user, got", count)
	}

	user, ok := testUsers.Get(UUID)
	if !ok || user.UUID != UUID {
		t.Fatal("Failed to retrieve UUID", UUID, "in users")
	}
}

func TestOnClose(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()

	msg := hook.OnConnection(UUID)
	msg = hook.OnClose(UUID)
	if msg == nil {
		t.Fatal("Expecting non nil message")
	}
	if msg.Event != websocket.USERDEL {
		t.Fatal("Expecting event", websocket.USERDEL, "got", msg.Event)
	}

	typedMessage, ok := msg.Data.(uuid.UUID)
	if !ok {
		t.Fatal("Expecting User structure, can't cast it:", msg.Data)
	}
	if typedMessage != UUID {
		t.Fatal("Expecting UUID", UUID, "got", typedMessage)
	}

	if count := testUsers.Count(); count != 0 {
		t.Fatal("Expecting 0 user, got", count)
	}
}

func TestOnSecondClose(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()

	msg := hook.OnConnection(UUID)
	msg = hook.OnClose(UUID)
	msg = hook.OnClose(UUID)
	if msg != nil {
		t.Fatal("Expecting nil message on second connection")
	}

	if count := testUsers.Count(); count != 0 {
		t.Fatal("Expecting 0 user, got", count)
	}
}

func TestOnAddEvent(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()
	name := "test"

	testMsg := websocket.Message{
		Event: websocket.USERADD,
		Data: users.User{
			UUID: UUID,
			Name: name,
		},
	}

	msg := hook.OnEvent(&testMsg)
	if msg == nil {
		t.Fatal("Expecting non nil message")
	}

	if count := testUsers.Count(); count != 1 {
		t.Fatal("Expecting 1 user, got", count)
	}

	user, ok := testUsers.Get(UUID)
	if !ok || user.UUID != UUID || user.Name != name {
		t.Fatal("Failed to retrieve UUID", UUID, "in users correctly")
	}
}

func TestOnSecondAddEvent(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()
	name := "test"

	testMsg := websocket.Message{
		Event: websocket.USERADD,
		Data: users.User{
			UUID: UUID,
			Name: name,
		},
	}

	msg := hook.OnEvent(&testMsg)
	msg = hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message")
	}

	if count := testUsers.Count(); count != 1 {
		t.Fatal("Expecting 1 user, got", count)
	}

	user, ok := testUsers.Get(UUID)
	if !ok || user.UUID != UUID || user.Name != name {
		t.Fatal("Failed to retrieve UUID", UUID, "in users correctly")
	}
}

func TestOnUpdateEvent(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()
	name := "test"

	testMsg := websocket.Message{
		Event: websocket.USERUPDATE,
		Data: users.User{
			UUID: UUID,
			Name: name,
		},
	}

	msg := hook.OnConnection(UUID)
	msg = hook.OnEvent(&testMsg)
	if msg == nil {
		t.Fatal("Expecting non nil message")
	}

	if count := testUsers.Count(); count != 1 {
		t.Fatal("Expecting 1 user, got", count)
	}

	user, ok := testUsers.Get(UUID)
	if !ok || user.UUID != UUID || user.Name != name {
		t.Fatal("Failed to retrieve UUID", UUID, "in users correctly")
	}
}

func TestOnSecondUpdateEvent(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()
	name := "test"

	testMsg := websocket.Message{
		Event: websocket.USERUPDATE,
		Data: users.User{
			UUID: UUID,
			Name: name,
		},
	}

	msg := hook.OnConnection(UUID)
	msg = hook.OnEvent(&testMsg)
	msg = hook.OnEvent(&testMsg)
	if msg == nil {
		t.Fatal("Expecting non nil message")
	}

	if count := testUsers.Count(); count != 1 {
		t.Fatal("Expecting 1 user, got", count)
	}

	user, ok := testUsers.Get(UUID)
	if !ok || user.UUID != UUID || user.Name != name {
		t.Fatal("Failed to retrieve UUID", UUID, "in users correctly")
	}
}

func TestOnDelEvent(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()
	name := "test"

	testMsg := websocket.Message{
		Event: websocket.USERDEL,
		Data: users.User{
			UUID: UUID,
			Name: name,
		},
	}

	msg := hook.OnConnection(UUID)
	msg = hook.OnEvent(&testMsg)
	if msg == nil {
		t.Fatal("Expecting non nil message")
	}

	if count := testUsers.Count(); count != 0 {
		t.Fatal("Expecting 0 user, got", count)
	}
}

func TestOnSecondDelEvent(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)
	UUID := uuid.New()
	name := "test"

	testMsg := websocket.Message{
		Event: websocket.USERDEL,
		Data: users.User{
			UUID: UUID,
			Name: name,
		},
	}

	msg := hook.OnConnection(UUID)
	msg = hook.OnEvent(&testMsg)
	msg = hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting non nil message")
	}

	if count := testUsers.Count(); count != 0 {
		t.Fatal("Expecting 0 user, got", count)
	}
}

func TestUnknownEvent(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)

	testMsg := websocket.Message{
		Event: "unknown",
		Data: users.User{
			UUID: uuid.New(),
			Name: "test",
		},
	}

	msg := hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message on second connection")
	}
}

func TestMalformedMessage(t *testing.T) {
	testUsers := users.Create()
	hook := CreateUserHook(&testUsers)

	testMsg := websocket.Message{
		Event: websocket.USERDEL,
		Data:  "malformed data",
	}

	msg := hook.OnEvent(&testMsg)
	if msg != nil {
		t.Fatal("Expecting nil message on second connection")
	}
}
