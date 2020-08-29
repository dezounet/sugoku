package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dezounet/sugokud/internal/websocket"
	"github.com/dezounet/sugokud/pkg/sudoku"
	"github.com/go-redis/redis/v8"
)

// GetGridHandler to serve HTTP GET request on sudoku grid
type GetGridHandler struct {
	Grid *sudoku.Grid
}

func (handler *GetGridHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setHeader(w.Header())

	json, err := json.Marshal(handler.Grid)
	if err != nil {
		log.Println("Failed getting grid: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	} else {
		log.Println("serving grid...")
		w.Write(json)
	}
}

// GetGridUUIDHandler to serve HTTP GET request on sudoku grid UUID
type GetGridUUIDHandler struct {
	Grid *sudoku.Grid
}

func (handler *GetGridUUIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setHeader(w.Header())

	json, err := json.Marshal(handler.Grid.UUID)
	if err != nil {
		log.Println("Failed getting grid UUID: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	} else {
		log.Println("serving grid uuid...")
		w.Write(json)
	}
}

// GetGridResetHandler to serve HTTP GET request to reset sudoku grid
type GetGridResetHandler struct {
	Grid      *sudoku.Grid
	Broadcast chan websocket.Message
	Redis     *redis.Client
}

func (handler *GetGridResetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setHeader(w.Header())

	if handler.Grid.IsSolved() {
		difficultyParams := r.URL.Query().Get("difficulty")

		difficulty := sudoku.EASY
		if difficultyParams == "medium" {
			difficulty = sudoku.MEDIUM
		} else if difficultyParams == "hard" {
			difficulty = sudoku.HARD
		} else if difficultyParams == "nightmare" {
			difficulty = sudoku.NIGHTMARE
		}

		// Initialize a new Grid
		handler.Grid.Initialize(difficulty)

		// Increment solved grid counter
		ok := sudoku.IncrementCounter(handler.Redis)
		if !ok {
			log.Println(
				"Failed to increment solved grid counter,",
				"this solved grid will not be added to total count",
			)
		}

		// Send response
		json, err := json.Marshal("OK")
		if err != nil {
			log.Println("Failed getting grid UUID: ", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		} else {
			log.Println("resetting grid...")
			w.Write(json)

			if handler.Broadcast != nil {
				msg := websocket.Message{
					Event: websocket.GRIDRESET,
				}
				handler.Broadcast <- msg
			}
		}
	} else {
		http.Error(w, http.StatusText(http.StatusForbidden),
			http.StatusForbidden)
	}
}

// GetGridCounter to serve HTTP GET request to get the number of sudoku grid solved
type GetGridCounter struct {
	Redis *redis.Client
}

func (handler *GetGridCounter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	setHeader(w.Header())

	count := sudoku.GetCounter(handler.Redis)
	json, err := json.Marshal(count)
	if err != nil {
		log.Println("Failed getting grid UUID: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError),
			http.StatusInternalServerError)
	} else {
		w.Write(json)
	}
}
