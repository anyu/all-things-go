// Fan-out/Fan-in

// Fan-out: Divide data into smaller chunks and distribute amongst works in parallel.
// Fan-in: Process multiple input channels and combine into a single channel.

// Example from: https://blog.golang.org/pipelines

package main

import (
	"fmt"
	"sync"
)

func main() {
	in := gen(2, 3)

	// Distribute the sq work across two goroutines that both read from in
	c1 := sq(in)
	c2 := sq(in)

	// Consume merged output
	for n := range merge(c1, c2) {
		fmt.Println(n) // 4 then 9, or 9 then 4
	}
}

// gen is a 'generator' that converts integers to channels
func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// sq operates on each integer channel
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// Start an output goroutine for each input channel.
	// Copy values from each input channel to out until input channel is closed.
	output := func(c <-chan int) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are done.
	// This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}