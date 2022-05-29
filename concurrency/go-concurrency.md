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

## Synchronizing goroutines

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

**Buffered** channels accept a limited number of values (its capacity) without needing a corresponding receive.

In other words, a buffered receive channel will block only if the buffer is empty.

A buffered send channel will block only if the buffer is full.

> Note: An unbuffered channel really is just a buffered channel with a capacity of 0.

> Use buffered channels with care as they hide possible deadlocks

```go
func main() {
  // removing the 2 here would cause a deadlock 
  // unless sending to dataStream is in a separate go routine
  dataStream := make(chan string, 2) 

  dataStream <- "msg 1"
  dataStream <- "msg 2"

  fmt.Println(<-dataStream)
  fmt.Println(<-dataStream)
}
```

### Closing channels

You can close a channel to signal that no more values will be sent to the channel.
Note that you can still *read* from a close channel, but you can't send anything more to it.

```go
dataStream := make(chan int)
close(dataStream)
val, ok := <-dataStream
fmt.Printf("%v: %v", ok, val) // false: 0
```

This can facilitate unblocking multiple goroutines waiting on the same channel (instead of writing `n` times to the channel to unblock each)

```go
func main() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)

	// Because we close'd the channel, the iteration ends after receiving 2 elements.
  // Removing the close line would result in a deadlock
	for e := range queue {
		fmt.Println(e)
	}
}
```

Closing or sending values to a closed channel panics.

### Nil channels

- Reading or writing to a `nil` channel will block (may deadlock)
  ```go
  var dataStream chan interface{}
  <-dataStream
  // fatal error: all goroutines are asleep - deadlock!
  ```
- Closing a `nil` channel will panic

### Channel ownership/handling

The goroutine that owns a channel should:
1. Instantiate the channel (avoids risk of writing to or closing a `nil` channel)
1. Perform writes or pass ownership to another goroutine
1. Close the channel (avoids risk of writing to a closed channel)
4. Encapsulate the above and expose them via a read channel

As a channel consumer, you should only have to worry about 2 things:
- knowing when a channel is closed // check 2nd return value from read operation
- handling the fact that reads can and will block

Keeping channel ownership scope small will make things easier to reason about.

```go
// Notice how resultStream's lifecycle is encapsulated within chanOwner
chanOwner := func() <-chan int {
  resultStream := make(chan int, 2)
  go func() {
    defer close(resultStream) // ensure channel is closed when done
    for i := 0; i <= 2; i++ {
      resultStream <- i
    }
  }()
  return resultStream // will be implicity converted to read-only due to the return type
}

resultStream := chanOwner()
for result := range resultStream { // range over channel to unblock
  fmt.Printf("Received: %d\n", result)
}
fmt.Println("Done receiving")
/*
Received: 0
Received: 1
Received: 2
Done receiving
*/
```

## Select

Selects help you compose channels together.

```go
var c1, c2 <- chan interface{}
var c3 chan<- interface{}

select {
case <- c1:
  // do something
case <- c2:
  // do something
case c3<- struct{}{}: 
  // do something  
}
```

While they look like `switch` blocks, these case statements aren't tested sequentially and execution does not 'fall through'. 

All channel read/writes are evaluted simultaneously. If no channels are ready, the entire `select` blocks. When one channel is ready, its statement executes. 

If no channel is ready, but you want to still take an action in the meantime, you can use a `default` case.

An empty select statement blocks forever:

```go
select {}
```

If there's nothing to do while channels are blocked, but you don't want to block forever, add a time out:

```go
var c<-chan int
select {
  case <-c: // blocks forever since we're reading from a nil channel
  case <- time.After (1*time.Second):
    fmt.Println("Timed out)
}
```
