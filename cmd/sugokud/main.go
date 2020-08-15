package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dezounet/sugokud/internal/api"
	"github.com/dezounet/sugokud/internal/users"
	"github.com/dezounet/sugokud/internal/websocket"
	hook "github.com/dezounet/sugokud/internal/websocket/hooks"
)

func main() {
	port := flag.Int("p", 8080, "Server is going to listen on this TCP port")
	flag.Parse()

	// Initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	// Create internal storage struct
	users := users.Create()

	// Configure API
	// - websocket
	websocketHandler := websocket.GetHandler()
	websocket.RegisterHook(&websocketHandler, hook.CreateUserHook(&users))
	http.Handle("/websocket", &websocketHandler)
	// - HTTP users
	userHandler := api.UserHandler{Users: &users}
	http.Handle("/users", &userHandler)

	// Start listening for incoming messages
	go websocket.HandleMessages(&websocketHandler)

	// Start server
	log.Println("http server started on :" + strconv.Itoa(*port))
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("Fail to listen and serve: ", err)
	}
}
