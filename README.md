# go-orchestrator

A very minimalistic orchestrator based on the Docker runtime written in Go. Courtesy of Tim Boring's book "Build An Orchestrator In Go"

## The Components of An Orchestration System

### Task

Tasks are the backbone component of an orchestrator. They are the **smallest unit of work** in an orchestration system.

**They can be thought of as containerized processes that run on a single machine (e.g., a Worker) inside a Container.**
For example, one Task would be to run NGINX, Rest API, or a microservice inside a Worker/Node.

Technically Tasks aren't exactly like Containerized Processes but abstractions created by the Orchestrator
to represent Containerized Processes. Why? Because **they include metadata such as the state, resource limits,
restart policies, and information on how they should be managed, monitored, and run by the Manager.**

As this Orchestrator uses Docker as its Container Runtime, Tasks run as Docker Containers.

### Worker

Workers are the muscles of an Orchestrator. **Tasks are assigned by the Manager to run as containers inside Workers.**

If a Task fails, Workers must attempt to restart it.

**Workers also make metrics and statistics about its Tasks available to the Manager** for the purpose of Scheduling.

### Job

### Scheduler

### Manager

### Cluster

### In Kubernetes' Terms

1. [Manager](#manager) = Control Plane
2. [Worker](#worker) = Kubelet/Node
3. [Job](#job) + [Task](#task) = [Kubernetes Job Objects](https://kubernetes.io/docs/concepts/workloads/controllers/job/)

The [Scheduler](#scheduler) and [Cluster](#cluster) are idiomatic and identical for Kubernetes without any differences.

## Go Learnings

Insights I've learned related to either Go or programming in general.

### Iota

Iota seems to be used in const group declarations (i.e., `const ()`), where its value increases after each line except empty or comment lines.

```go
const (
  Pending   State = iota // iota = 0
  Scheduled              // iota = 1
  Running                // iota = 2
  Completed              // iota = 3
  Failed                 // iota = 4
)
```

See more at the Go wiki: [Iota](https://go.dev/wiki/Iota).

### In-Memory Databases (IMDBs)

In Go we can create a quick and simple key-value In-Memory Database by using the `Map` type. Map declarations usually look like this `map[KeyType]ValueType`.

So, we could create a simple In-Memory Database for the Worker struct:

```go
type Worker struct {
  ...
  Db        map[uuid.UUID]*task.Task
}
```

The keys will be of type `uuid.UUID` and the values of type `*task.Task` (*i.e., pointers to Tasks*).

**For a more sophisticated IMDB, consider [bbolt](https://github.com/etcd-io/bbolt), which is an embedded key-value Database for Go with features such as Disk Persistence.**

See more at Wikipedia: [In-Memory Database](https://en.wikipedia.org/wiki/In-memory_database)
