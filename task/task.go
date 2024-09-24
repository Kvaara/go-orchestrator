package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
)

type Task struct {
	ID            uuid.UUID
	ContainerID   string
	Name          string
	State         State
	Image         string
	Cpu           float64
	Memory        int64
	Disk          int64
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

// Configuration for our Tasks.
type Config struct {
	Name          string // Identifies a task in our Orchestration system
	AttachStdin   bool
	AttachStdout  bool
	AttachStderr  bool
	ExposedPorts  nat.PortSet
	Cmd           []string
	Image         string // Specifies the Docker image of the container
	Cpu           float64
	Memory        int64    // Used by the Scheduler to find a node capable of running a Task
	Disk          int64    // Used by the Scheduler to find a node capable of running a Task
	Env           []string // Utilized to inject environment variables into containers.
	RestartPolicy string   // empty string, `always`, `unless-stopped`, or `on-failure`.
}

// Helper function for returning a copy of a Task's Config.
func NewConfig(t *Task) *Config {
	return &Config{
		Name:          t.Name,
		ExposedPorts:  t.ExposedPorts,
		Image:         t.Image,
		Cpu:           t.Cpu,
		Memory:        t.Memory,
		Disk:          t.Disk,
		RestartPolicy: t.RestartPolicy,
	}
}

// A helper function for returning a new Docker management object with a new desired Config.
func NewDocker(c *Config) *Docker {
	dc, _ := client.NewClientWithOpts(client.FromEnv)
	return &Docker{
		Client: dc,
		Config: *c,
	}
}

type Docker struct {
	Client *client.Client // Client object will be used to interact with Docker API
	Config Config
}

// A wrapper for aggregating standard information from methods that start/run containers.
type DockerResult struct {
	Error       error
	Action      string
	ContainerId string
	Result      string
}

/*
Runs a container by the following process:
 1. Pulls a desired image
 2. Configuring the container (e.g., settings a Restart Policy)
 3. Creates and starts the container
 4. Prints its logs to the terminal
 5. Returns a wrapped result.
*/
func (d *Docker) Run() DockerResult {
	// Context is a Go package that is utilized to manage deadlines, cancellation signals,
	// and other request-scoped values across API boundaries and between processes.
	ctx := context.Background()

	// Pulls the desired image from Dockerhub
	reader, err := d.Client.ImagePull(
		ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	// Copies `reader` to stdout (src to dst). In other words, prints `reader` to standard output for transparency/debugging
	// Copies data until `io.EOF` is reached
	io.Copy(os.Stdout, reader)

	// Configures the container's restart policy using our Task `Config` struct
	rp := container.RestartPolicy{
		Name: container.RestartPolicyMode(d.Config.RestartPolicy),
	}

	// Holds the resources required by the container. Memory is specified via our Tasks's `Config` struct.
	r := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.Cpu * math.Pow(10, 9)),
	}

	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	// Configuration of the host machine where the container runs (e.g., Linux machine)
	hc := container.HostConfig{
		RestartPolicy: rp,
		Resources:     r,
		// Docker will expose all ports automatically by randomly choosing available ports on the host.
		PublishAllPorts: true, // Same as passing -P to `docker run` which will publish all `EXPOSE`'d ports to a random port.
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	// Start the created container with empty `StartOptions` using its ID
	if err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		log.Printf("Error starting container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	// We need to get the logs from our container to print them to stdout for transparency/debugging purposes
	out, err := d.Client.ContainerLogs(
		ctx,
		resp.ID,
		container.LogsOptions{ShowStdout: true, ShowStderr: true},
	)
	if err != nil {
		log.Printf("Error getting logs for container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	// Same as `io.copy` but can be used with streams that contain both stdout and stderr information.
	// Useful for Docker container logs because the output includes both stdout and stderr in a single stream, and
	// needs to be "demultiplexed".
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// Return `DockerResult` with the Container ID wrapped
	return DockerResult{ContainerId: resp.ID, Action: "start", Result: "success"}
}

// Stops and removes a container using its ID.
func (d *Docker) StopAndRemove(id string) DockerResult {
	log.Printf("Attempting to stop container %v", id)
	ctx := context.Background()

	// Try to stop the container using its ID:
	err := d.Client.ContainerStop(ctx, id, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}

	// If stopping is successful, try to kill and remove the container using its ID from the host machine:
	err = d.Client.ContainerRemove(ctx, id, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
		Force:         false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", id, err)
		return DockerResult{Error: err}
	}

	return DockerResult{Action: "stop", Result: "success", Error: nil}
}
