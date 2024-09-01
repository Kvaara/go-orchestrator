/*
Tasks are the backbone component of an orchestrator. They are the smallest unit of work in an orchestration system.

They can be thought of as containerized processes that run on a single machine (e.g., a Worker) inside a Container.
For example, one Task would be to run NGINX, Rest API, or a microservice inside a Worker/Node.

Technically Tasks aren't exactly like Containerized Processes but abstractions created by the Orchestrator
to represent Containerized Processes. Why? Because they include metadata such as the state, resource limits,
restart policies, and information on how they should be managed, monitored, and run by the Manager.

As this Orchestrator uses Docker as its Container Runtime, Tasks run as Docker Containers.
*/

package task

import (
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type State int

const (
	Pending State = iota
	Scheduled
	Running
	Completed
	Failed
)

type Task struct {
	ID            uuid.UUID
	Name          string
	State         State
	Image         string
	Memory        int
	Disk          int
	ExposedPorts  nat.PortSet
	PortBindings  map[string]string
	RestartPolicy string
	StartTime     time.Time
	FinishTime    time.Time
}

type TaskEvent struct {
	ID        uuid.UUID
	State     State
	Timestamp time.Time
	Task      Task
}
