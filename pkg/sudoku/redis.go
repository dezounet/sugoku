package sudoku

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var keyGrid = "grid"
var keyCounter = "counter"

// Connect to a redis DB
func Connect(addr string, passwd string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   passwd,
		MaxRetries: 3,
	})

	return client
}

// SetGrid to redis instance, in order to persist it
func SetGrid(client *redis.Client, grid *Grid) bool {
	ok := false

	if client != nil {
		// Serializing grid to JSON
		json, err := json.Marshal(*grid)
		if err != nil {
			log.Println("Failed serialzing grid: ", err)
		} else {
			err = client.Set(ctx, keyGrid, json, 0).Err()
			if err != nil {
				log.Println("Unexpected error caching grid:", err)
			} else {
				ok = true
			}
		}
	}

	return ok
}

// GetGrid from redis instance, and return a Grid object
func GetGrid(client *redis.Client) *Grid {
	var grid *Grid = nil

	if client != nil {
		// Get from redis
		jsonGrid, err := client.Get(ctx, keyGrid).Result()
		if err == redis.Nil {
			log.Println("No cached grid found")
		} else if err != nil {
			log.Println("Unexpected error getting cached grid:", err)
		} else {
			// Unmarshall sotred data
			retrievedGrid := Grid{}
			err = json.Unmarshal([]byte(jsonGrid), &retrievedGrid)
			if err == nil {
				grid = &retrievedGrid
			} else {
				grid = nil
			}
		}
	}

	return grid
}

// IncrementCounter of solved sudoku grid in redis
func IncrementCounter(client *redis.Client) bool {
	ok := false

	if client != nil {
		err := client.Incr(ctx, keyCounter).Err()
		if err != nil {
			log.Println("Unexpected error caching counter:", err)
		} else {
			ok = true
		}
	}

	return ok
}

// GetCounter of total solved sudoku grid from redis
func GetCounter(client *redis.Client) int {
	count := 0

	if client != nil {
		var err error
		count, err = client.Get(ctx, keyCounter).Int()
		if err != nil {
			log.Println("Unexpected error getting solved grid counter:", err)
			count = 0
		}
	}

	return count
}
