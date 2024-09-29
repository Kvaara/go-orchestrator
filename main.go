package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/kvaara/go-orchestrator/task"
	"github.com/kvaara/go-orchestrator/worker"
)

func main() {
	host := os.Getenv("CUBE_HOST")
	// Converts string to int:
	port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))

	fmt.Println("Starting Cube worker")

	w := worker.Worker{
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{Address: host, Port: port, Worker: &w}

	// The `go` keyword is used to create Goroutines (i.e., threads) to handle concurrent operations.
	// This means that the below loop is nonblocking and moves on to the next statement.
	go runTasks(&w)

	go api.ServeAPI()

	// This is a common idiom to insert blocks. A select blocks until one of its `case`s can be run.
	// An empty select block works as an eternal block ensuring that the main function never returns so the Go application never stops
	select {}
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			result := w.RunTaskInQueue()
			if result.Error != nil {
				log.Printf("Error running task: %v\n", result.Error)
			}
		} else {
			log.Printf("No tasks to process currently.\n")
		}
		log.Println("Sleeping for 10 seconds.")
		time.Sleep(10 * time.Second)
	}
}
