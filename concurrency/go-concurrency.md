# Go Concurrency

## The Basics

The first thing to clear up is the difference between **concurrency** and **parallelism**, which are often used interchangeably, but is not the same.

- **Concurrency** is about juggling multiple tasks at once in a window of time (execution's start and end times overlap), but not actually executing them at the same time. 

- **Parallelism** is actually executing tasks at the same time.

An analogy helps here. Doing something concurrently is like making a meal - you kick off boiling pasta while chopping onions. Parallelism would be having helpers chopping onions simultaneously.

Parallel tasks must be executed on different hardware (separate CPU cores); concurrent tasks could be executed on the same hardware.

## Understanding OS processes, threads

Before we get into Go specific terms, let's get acquainted with some general concepts.

Most people have a clear idea of what a **computer program** is, so we'll start there. You have Slack, Google Chrome, Evernote, etc. They are stored on disk when not running.

Well, when a program starts, it needs OS resources to run (eg. memory)

- A **process** is essentially a program that's been loaded into memory along with the resouces it needs. The role of an OS, is to allocate resources to various processes.

  Processes are isolated and independent of each other; they do not share resources. Switching processes is slow.

- A **thread** is 1 unit of execution within a process; a process can have many threads. Threads have their own contexts, but share certain resources (since the whole point is to reduce cost of context switching) allocated to the process. 

  They are sometimes called 'lightweight processes', because the cost of communication/switching between threads is low. But the tradeoffs is they can affect each other.

A single-threaded process is what it sounds like: a process that only has 1 thread.

In a multi-threaded process, threads have their own stacks, but share the heap. 

## Goroutines

In Go, the basic unit of organization is a **goroutine**. 

You can think of a goroutine as a lightweight thread, or a function that can run concurrently with other code.

Whereas threads are managed by the OS, however, goroutines are managed by Go's runtime. You can easily communicate between two goroutines, whereas you can't with threads. 

> Note: Green threads are managed by a language's runtime, so goroutines are sometimes considered green threads, but technically they are in the category of 'coroutines'. Skipping this detail for now.

Goroutines are also very cheap; they're initialized with only 2KB of stack memory and grows when necessary. Whereas in Java, a thread is allocated a fixed memory size.

Every Go program has at least 1 goroutine: the `main` goroutine, which is automatically created when the program runs.

### Creating a goroutine

Let's create a goroutine via the `go` keyword:

```go
func main() {
  go hello()
  // do something else
}

func hello() {
  fmt.Println("hello")
}
```

We can also create one via an anonymous function:

```go
func main() {
  go func() {
    fmt.Println("hello")
  }()
  // do something else
}
```

or assigning it to a variable:

```go
func main() {
  hello := func() {
    fmt.Println("hello")
  }
  go hello()
  // do something else
}
```

### Synchronizing goroutines

The above examples have a problem in that the `main` goroutine will exit before the other goroutine has a chance to run. The simplest way to coordinate/synchronize goroutines is via **WaitGroups**.

A **WaitGroup** blocks a program's execution until the goroutines in it have executed. You can think of it like using a counter: 
- `Add` increments the counter
- `Done` decrements it
- `Wait` blocks until the counter is 0

```go
func main() {

  var wg sync.WaitGroup

  wg.Add(1)

  go func() {
    defer wg.Done()
    fmt.Println("hello")
  }()

  wg.Wait() // blocks until Done() completes
  // do something else
}
```

## Communicating between goroutines

Often, we don't just need to kick off concurrent execution, but also need to act on the results.

**Channels** are the mechanism via which goroutines can communicate/share data bidirectionally with each other.

```go
func main() {
  // create a channel via `make(chan TYPE)`
  // one convention is to name a channel variable `stream`
  dataStream := make(chan string) 
  go func() {
    dataStream <- "ping" // send a value _to_ the channel
  }()

  result := <-dataStream // receive a value _from_ a channel
  fmt.Println(result)
}
```

By default, channels are bidirectional; you can send or receive data on it. You can constrain a channel to be unidirectional:

```go
// receive only
dataStream := make(<-chan string)

// send only
dataStream := make(chan<- string)
```

This is more common in function parameters and return types.

### Buffered channels

Channels block. 

- A goroutine that wants to write to a channel that's full will wait until the channel is empty.
- similarly, it'll wait for a channel to be empty before reading from it.

By default, channels are **unbuffered** - they'll only accept sends if there's a corresponding receive that's ready to receive the value.

**Buffered** channels accept a limited number of values without needing a corresponding receive.

In other words, a buffered receive channel will block only if there's no value in the channel to receive.

A buffered send channel will block only if there's no available buffer to place the value being sent.

```go
func main() {
  dataStream := make(chan string, 2) // removing the 2 here will cause a deadlock - TODO: explain

  dataStream <- "msg 1"
  dataStream <- "msg 2"

  fmt.Println(<-dataStream)
  fmt.Println(<-dataStream)
}
```

### Closing channels


## Handling errors

