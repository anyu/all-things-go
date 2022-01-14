### Switch vs Select statements

#### Comparison
`select` is only used with channels.
`switch` is used with concrete types.

#### Switch statements
- goes in sequence 
- only falls through if explicitly stated
```go
someVar := 1
switch someVar {
case 0:
    fmt.Println(0)
    fallthrough
case 1:
    fmt.Println(1)
    fallthrough
case 2:
    fmt.Println(2)
}
```
- can match multiple for each case
```go
switch someVar {
case "aaa", "bbb":
    fmt.Printf("category 1")
case "ccc", "ddd":
    fmt.Printf("category 2")
case "eee", "fff":
    fmt.Printf("category 3")
}
```


#### Select statements
- lets a goroutine wait on multiple communication operations
- blocks until one of its cases can run, then it executes that case. It chooses one randomly if multiple are ready.
- empty select blocks forever, results in deadlock, panics

```go
select {}
```

- used for cancelling contexts
```go
select {
case <-ctx.Done():
    err = ctx.Err()
    break
case time.After(TODO: some time):
}
https://github.com/nytimes/auth-docs/blob/main/docs/team/interviews/flakylib/README.md
```