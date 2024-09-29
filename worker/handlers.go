package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/kvaara/go-orchestrator/task"
)

// Handlers are functions that can respond to our requests.

// Handler function for starting the Task specified in the request body.
func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	te := task.TaskEvent{}
	// Unmarshals (i.e., decodes) the request body's data into a type of `task.TaskEvent` using decoder `d`.
	// Decoding (reading) is unmarshaling and encoding (writing) is marshaling.
	err := d.Decode(&te)
	if err != nil {
		msg := fmt.Sprintf("Error unmarshalling body: %v\n", err)
		log.Print(msg)
		w.WriteHeader(400)
		e := ErrResponse{
			HTTPStatusCode: 400,
			Message:        msg,
		}
		json.NewEncoder(w).Encode(e)
		return
	}

	a.Worker.AddTaskToQueue(te.Task)
	log.Printf("Added task %v\n", te.Task.ID)
	w.WriteHeader(201)
	// Marshals (i.e., encodes) the Task into the response's body by using the response writer `w`.
	json.NewEncoder(w).Encode(te.Task)
}

// Handler function for fetching all the Task's in the Worker's DB.
func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(a.Worker.GetTasksInDB())
	}
}

// Handler function for stopping the Task with the ID specified in the request path.
func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	taskID := r.PathValue("taskID")
	if taskID == "" {
		log.Printf("No taskID passed in request.\n")
		w.WriteHeader(400)
	}

	tID, _ := uuid.Parse(taskID)
	taskToStop, ok := a.Worker.Db[tID]
	if !ok {
		log.Printf("No task with ID %v found", tID)
		w.WriteHeader(404)
	}

	// we need to make a copy so we are not modifying the task in the datastore
	taskCopy := *taskToStop
	taskCopy.State = task.Completed
	a.Worker.AddTaskToQueue(taskCopy)

	log.Printf("Added task %v to stop container %v\n", taskToStop.ID, taskToStop.ContainerID)
	w.WriteHeader(204)
}
