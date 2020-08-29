package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dezounet/sugokud/internal/users"
)

// GetUsersHandler to serve HTTP GET request on users
type GetUsersHandler struct {
	Users *users.LockableUsers
}

func (handler *GetUsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setHeader(w.Header())

	samples := handler.Users.Sample(10)

	json, err := json.Marshal(samples)
	if err != nil {
		log.Println("Failed getting user samples: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	} else {
		log.Println("serving users...")
		w.Write(json)
	}
}
