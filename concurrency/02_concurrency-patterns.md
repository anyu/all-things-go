# Concurrency Patterns

## The for-select loop

Here's a pattern you'll commonly see:

```go
for { // either an infinite loop, or ranging over something
  select {
    // do some channel work
  }
}
```

Some scenarios for this include:

**1. Sending iteration values out on a channel**
  - when you want to convert something that can be iterated over into values for a channel
  ```go
  for _, s := range []string{"a", "b", "c"} {
    select {
    case <-done:
      return
    case dataStream <- s:
    }
  }
  ```

**2. Looping infinitely until stopped**
  -  variation 1 - if done isn't closed, exit select and execute rest of loop body:
  ```go
  for {
    select {
    case <-done:
      return
    default:
    }
    // do non-preemptable work
  }
  ```

  - varation 2 - if done isn't closed, execute default clause:
  ```go
  for {
    select {
    case <-done:
      return
    default:
      // do non-preemptable work
    }
  }
  ```

## Preventing Goroutine leaks

Goroutines terminate when:
- it completes its work
- it can't continue its work due to an unrecoverable error
- it's told to stop working

While goroutines are cheap, they still cost some resources (they're not garbage collected).

If a parent goroutine ends, we should make sure to clean up its children goroutines.

Conventionally, this is done via the `done` channel. 

```go
// pass done channel to function (done channel is first parameter by convention)
doWork := func(done <-chan interface{}, strings <-chan string) <-chan interface{} {
  terminated := make(chan interface{})
  go func() {
    defer fmt.Println("doWork exited.")
    defer close(terminated)
    for {
      select {
      case s := <-strings:
        // do something
        fmt.Println(s)
      case <-done: // return if done channel has received
        return
      }
    }
  }()
  return terminated
}

done := make(chan interface{})
terminated := doWork(done, nil)

go func() { // goroutine to cancel the doWork goroutine if it's been more than 1 sec
  // cancel operation after 1 sec
  time.Sleep(1 * time.Second)
  fmt.Println("Canceling doWork goroutine...")
  close(done)
}()
<-terminated // join doWork goroutine with main goroutine
fmt.Println("Done")
```

## Error Handling

Delegate error handling to the caller of the goroutine, who is in a better position to act on the error.

One approach is to return a struct that contains the possible outcomes from the goroutine, the desired data or error.

```go
type Result struct {
  Error    error
  Response *http.Response
}

checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {
  results := make(chan Result)
  go func() {
    defer close(results)

    for _, url := range urls {
      var result Result

      resp, err := http.Get(url)
      result = Result{
        Error:    err,
        Response: resp,
      }
      select {
      case <-done:
        return
      case results <- result:
      }
    }
  }()
  return results
}
done := make(chan interface{})
defer close(done)

urls := []string{"https://www.google.com", "https://badhost"}
for result := range checkStatus(done, urls...) {
  if result.Error != nil {
    fmt.Printf("error: %v\n", result.Error)
    continue
  }
  fmt.Printf("Response: %v\n", result.Response.Status)
}
```

## Pipelines

Pipelines facilitate transforming data. 

Each stage is a function that takes in and returns the same type, so they can be chained.

### Batch processing

Each stage receives/operates on/outputs chunks of data instead of 1 discrete element at a time.

Cons: Memory footprint is a little larger since we need a new slice of equal length of input data to store results of calculations.

### Stream processing

Each stage receives/operates on/outputs one element at a time.

Cons: Need to rejigger where the pipelinng happens; may limit ability to scale.










