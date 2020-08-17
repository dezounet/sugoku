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
	"github.com/dezounet/sugokud/pkg/sudoku"
)

func main() {
	port := flag.Int("p", 8080, "Server is going to listen on this TCP port")
	size := flag.Int("s", 3, "Block size to use for the generated grid")
	flag.Parse()

	// Initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	// Create internal storage struct
	users := users.Create()
	grid := sudoku.CreateEmptyGrid(*size)

	// Configure API
	// - websocket
	websocketHandler := websocket.GetHandler()
	websocket.RegisterHook(&websocketHandler, hook.CreateUserHook(&users))
	websocket.RegisterHook(&websocketHandler, hook.CreateGridHook(grid))
	http.Handle("/websocket", &websocketHandler)
	// - HTTP grid
	getGridHandler := api.GetGridHandler{Grid: grid}
	http.Handle("/grid", &getGridHandler)
	getUUIDHandler := api.GetGridUUIDHandler{Grid: grid}
	http.Handle("/grid/uuid", &getUUIDHandler)
	getGridResetHandler := api.GetGridResetHandler{
		Grid:      grid,
		Broadcast: websocketHandler.Broadcast,
	}
	http.Handle("/grid/reset", &getGridResetHandler)
	// - HTTP users
	userHandler := api.GetUsersHandler{Users: &users}
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
