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

You can think of a goroutine as a lightweight thread. 

Whereas threads are managed by the OS, however, goroutines are managed by Go's runtime. You can easily communicate between two goroutines, whereas you can't with threads. 

> Note: Green threads are managed by a language's runtime, so goroutines are sometimes considered green threads, but technically they are in the category of 'coroutines'. Skipping this detail for now.

Goroutines are also very cheap; they're initialized with only 2KB of stack memory and grows when necessary. Whereas in Java, a thread is allocated a fixed memory size.

Every Go program has at least 1 goroutine: the `main` goroutine, which is automatically created when the program runs.

## Kicking off an async operation

## Receiving data back from async operation

## Handling errors

