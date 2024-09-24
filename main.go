package main

import (
	"fmt"
	"time"

	"github.com/golang-collections/collections/queue"

	"github.com/google/uuid"
	"github.com/kvaara/go-orchestrator/task"
	"github.com/kvaara/go-orchestrator/worker"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)
	w := worker.Worker{
		Name:      "Worker-1",
		Queue:     *queue.New(),
		Db:        db,
		TaskCount: 0,
	}

	fmt.Printf("Queue should be empty so no Tasks should be run")
	result := w.RunTaskInQueue()
	if result.Error != nil {
		panic(result.Error)
	}
	time.Sleep(time.Second * 10)

	// Hmmm... The way we create Tasks isn't ideal.
	// It's hard to tell what properties should be defined (i.e., what properties are required).
	t := task.Task{
		ID:    uuid.New(),
		Name:  "Task-1",
		State: task.Scheduled,
		Image: "strm/helloworld-http",
	}

	fmt.Println("Queue should have one Scheduled Task, which should be dequeued, run, and changed to Running.")
	w.AddTaskToQueue(t)
	result = w.RunTaskInQueue()
	if result.Error != nil {
		panic(result.Error)
	}
	t.ContainerID = result.ContainerId

	fmt.Printf("task %s is running in container %s\n", t.ID, t.ContainerID)
	time.Sleep(time.Second * 30)

	fmt.Printf("Queue should be empty so we add the same Task stated as Completed to stop it.")
	t.State = task.Completed
	w.AddTaskToQueue(t)
	time.Sleep(time.Second * 30)

	fmt.Printf("Queue should have one Completed Task, which should be dequeued, stopped, and changed to Completed: %s", t.ID)
	result = w.RunTaskInQueue()
	if result.Error != nil {
		panic(result.Error)
	}
}
