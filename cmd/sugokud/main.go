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
	"github.com/go-redis/redis/v8"
)

func periodicBackup(redis *redis.Client, grid *sudoku.Grid) {
	t := time.NewTicker(time.Minute)
	for {
		log.Println("Updating cache...")
		ok := sudoku.SetGrid(redis, grid)
		if !ok {
			log.Println("Failed storing new grid in cache, if you restart any data will be lost")
		}

		// Wait
		<-t.C
	}
}

func main() {
	port := flag.Int("p", 8080, "Server is going to listen on this TCP port")
	size := flag.Int("s", 3, "Block size to use for the generated grid")
	redisAddr := flag.String("redis", "redis:6379", "Redis DB address & port")
	redisPasswd := flag.String("passwd", "ThisIsNotSecure", "Password used to connect to redis")

	flag.Parse()

	// Initialize global pseudo random generator
	rand.Seed(time.Now().Unix())

	// Create internal storage struct
	users := users.Create()

	// Connect to redis and get grid
	redis := sudoku.Connect(*redisAddr, *redisPasswd)
	grid := sudoku.GetGrid(redis)
	if grid == nil {
		// No grid found in cache, creating a new one
		grid = sudoku.CreateEmptyGrid(*size)
		grid.Initialize(sudoku.EASY)

		// Storing this new cached grid
		ok := sudoku.SetGrid(redis, grid)
		if !ok {
			log.Println("Failed storing new grid in cache, if you restart any data will be lost")
		}
	}

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
		Redis:     redis,
	}
	http.Handle("/grid/reset", &getGridResetHandler)
	getGridCounterHandler := api.GetGridCounter{Redis: redis}
	http.Handle("/grid/count", &getGridCounterHandler)

	// - HTTP users
	userHandler := api.GetUsersHandler{Users: &users}
	http.Handle("/users", &userHandler)

	// Start listening for incoming messages
	go websocket.HandleMessages(&websocketHandler)

	// Background cache update
	go periodicBackup(redis, grid)

	// Start server
	log.Println("http server started on :" + strconv.Itoa(*port))
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("Fail to listen and serve: ", err)
	}
}
