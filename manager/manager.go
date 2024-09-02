package manager

import (
	"fmt"

	"github.com/kvaara/go-orchestrator/task"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
)

type Manager struct {
	Pending queue.Queue
	TaskDb  map[string][]*task.Task
	EventDb map[string][]*task.TaskEvent
	Workers []string
	// A Map of Workers to a list of their Tasks
	WorkerTaskMap map[string][]uuid.UUID
	// A Map of Tasks to the Workers under where they run
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {
	fmt.Println("I will select an appropriate worker")
}

func (m *Manager) UpdateTasks() {
	fmt.Println("I will update tasks")
}

func (m *Manager) SendWork() {
	fmt.Println("I will send work to workers")
}
