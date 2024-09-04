# go-orchestrator

A very minimalistic orchestrator based on the Docker runtime written in Go. Courtesy of Tim Boring's book "Build An Orchestrator In Go"

Below is an overall architecture of the orchestrator:
![Overall Architecture of the Orchestrator](/attachments/overall-architecture.png)

## The Components of An Orchestration System

### Task

Tasks are the backbone component of an orchestrator. They are the **smallest unit of work** in an orchestration system.

**They can be thought of as containerized processes that run on a single machine (e.g., a Worker) inside a Container.**
For example, one Task would be to run NGINX, Rest API, or a microservice inside a Worker/Node.

Technically Tasks aren't exactly like Containerized Processes but abstractions created by the Orchestrator
to represent Containerized Processes. Why? Because **they include metadata such as the state, resource limits,
restart policies, and information on how they should be managed, monitored, and run by the Manager.**

As this Orchestrator uses Docker as its Container Runtime, Tasks run as Docker Containers.

![Process of Task States](/attachments/task-states.png)

### Worker

Workers are the muscles of an Orchestrator. **Tasks are assigned by the Manager to run as containers inside Workers.**

If a Task fails, Workers must attempt to restart it.

**Workers also make metrics and statistics about its Tasks available to the Manager** for the purpose of Scheduling.

Workers deal with the logical workload (the [Tasks](#task)) of the Orchestrator. Workers are a type of [Node](#node)

### Node

A node is an object that represents any machine in the [Cluster](#cluster). Types of nodes, for example:

- The [Manager](#manager)
- A [Worker](#worker)

While Workers handle the logical workload, **Nodes are either physical or virtual machines.**

### Job

### Scheduler

The Scheduler is an advisor to the [Manager](#manager) providing it the following information:

1. Determines a set of candidate [Workers](#worker) on which a [Task](#task) could run.
2. **Scores the preceding candidates by using a scheduling algorithm.**
3. Picks the candidate with the best score.

There can be different kinds of schedulers, which is the main reason why **the Scheduler is implemented as an Interface type.**

### Manager

The Manager is the brains of the Orchestrator. It **accepts requests coming from the user** (e.g., via the [CLI](#command-line-interface-cli-tool)) to start/stop [Tasks](#task).

The Manager **uses the [Scheduler](#scheduler) as an advisor in determining the best [Worker](#worker) candidate to whom to schedule a Task to be run.**

The Manager **puts Tasks into a FIFO (*First In First Out*) queue**.

The Manager also **collects metrics and statistics from the Workers**, which are then utilized by the [Scheduler](#scheduler).

The Manager will need to keep track of the Workers in the [Cluster](#cluster).

### Cluster

### Command Line Interface (CLI) Tool

### In Kubernetes' Terms

1. [Manager](#manager) = Control Plane
2. [Worker](#worker) = Kubelet/Node
3. [Job](#job) + [Task](#task) = [Kubernetes Job Objects](https://kubernetes.io/docs/concepts/workloads/controllers/job/)

The [Scheduler](#scheduler) and [Cluster](#cluster) are idiomatic and identical for Kubernetes without any differences.

## Go (and Other) Lessons Learned

Insights I've learned related to either Go or programming in general.

### UUID

UUIDs are Universally Unique IDentifiers. **They are 128-bits and, in practice, unique.** They follow a specific structure and a set of rules defined by [RFC 9562](https://datatracker.ietf.org/doc/html/rfc9562), which displaced the previous RFC 4122.

In theory, however, there's the possibility of two identical UUIDs being generated but that probability is extremely low.

### Interfaces

**Interfaces in Go support [Polymorphism](https://www.techtarget.com/whatis/definition/polymorphism).** This means that a type, which implements the interface type, can be used wherever the interface type is expected.

**They also define methods a type must implement to be considered an interface type.**

See more at Effective Go: [Interfaces and Other Types](https://go.dev/doc/effective_go#interfaces_and_types)

### In-Memory Databases (IMDBs)

In Go we can create a quick and simple key-value In-Memory Database by using the `Map` type. Map declarations usually look like this `map[KeyType]ValueType`.

- See [Go Maps In Action](https://go.dev/blog/maps).

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

### `make`

The `make` built-in function seems to be very useful in Go. It can return an initialized object of type Map, Slice, or Channel.

- Make is useful for data structures that require Runtime Initialization.

Unlike `new`, `make`'s return type is EXACTLY the same as the type of its argument, not a pointer to it. Also, `new` returns a zeroed value of a given type, which is useful for data types like `struct`s.

For example:

```go
  m := manager.Manager{
    ...
    TaskDb:  make(map[string][]*task.Task),
    EventDb: make(map[string][]*task.TaskEvent),
  }
```

Read more: [The new() vs make() Functions in Go – When to Use Each One](https://www.freecodecamp.org/news/new-vs-make-functions-in-go/)

### Pointers and Dereferences

- Pointers (e.g., `*Queue`) are values to memory addresses.
- Dereferences (e.g. `**Queue`) expose the actual value of the Pointer's memory address (*i.e., the data*).
- Ampersands (e.g., `&x`, where `x := 10`) are used to get the memory address of a variable. In other words, they create Pointers.
