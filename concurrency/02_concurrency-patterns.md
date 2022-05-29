# Concurrency Patterns

## Confinement

TODO

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
