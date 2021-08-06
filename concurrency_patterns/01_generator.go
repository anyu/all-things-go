// Generator (aka. Iterator)

// A function that returns a channel.

// Good for generating a sequence of values used to produce an output,
// as it allows the consumer to run in parallel while producer is processing next value.

// Example from: https://blog.golang.org/pipelines

package main

import "fmt"

func main() {
	ch := gen(0, 1, 2, 3) // function that returns a channel
	for n := range ch {
		fmt.Println(n)
	}
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out) // close channel when done to prevent deadlock
	}()
	return out
}

// Fibonacci example
func main() {
	for i := range fib(100) {
		fmt.Println(i)
	}
}

func fib(n int) chan int {
	out := make(chan int)
	go func() {
		for i, j := 0, 1; i < n; i, j = i+j, i {
			out <- i
		}
		close(out) // close channel when done to prevent deadlock
	}()
	return out
}
