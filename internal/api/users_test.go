package api

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/dezounet/sugokud/internal/users"
	"github.com/google/uuid"
)

func TestGetEmptyUsers(t *testing.T) {
	testUsers := users.Create()
	handler := GetUsersHandler{
		Users: &testUsers,
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, nil)

	data := make([]users.User, 0)
	err := json.Unmarshal(rec.Body.Bytes(), &data)
	if err != nil {
		t.Fatal("Failed to understand response:", rec.Body.String())
	}

	if len(data) > 0 {
		t.Fatal("Expecting empty users, got", data)
	}
}

func TestGetUsers(t *testing.T) {
	testUsers := users.Create()
	handler := GetUsersHandler{
		Users: &testUsers,
	}

	uuid := uuid.New().String()
	name := "test"
	user := users.User{
		UUID: uuid,
		Name: name,
	}

	testUsers.AddOrUpdate(user)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, nil)

	data := make([]users.User, 0)
	err := json.Unmarshal(rec.Body.Bytes(), &data)
	if err != nil {
		t.Fatal("Failed to understand response:", rec.Body.String())
	}

	if len(data) != 1 {
		t.Fatal("Expecting 1 user, got", data)
	}

	if data[0].UUID != uuid || data[0].Name != name {
		t.Fatal("Expecting user ", uuid, name, ", got", data[0])

	}
}
