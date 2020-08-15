package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dezounet/sugokud/internal/users"
)

// UserHandler to serve HTTP request on users
type UserHandler struct {
	Users *users.LockableUsers
}

func (handler *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	samples := handler.Users.Sample(10)

	json, err := json.Marshal(samples)
	if err != nil {
		log.Fatal("Failed getting user samples: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	} else {
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}
