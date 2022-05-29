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
