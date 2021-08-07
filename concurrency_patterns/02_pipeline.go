// Pipeline

// A series of stages connected by channels, where each stage is a group of goroutines running the same function.

// In each intermediary stage, the goroutines:
// - receive values (via inbound channels) -> perform operation -> send modified values out (via outbound channels)
//
// The first stage (aka. producer/source) only has outbound channels.
// The last stage (aka. consumer/sink) only has inbound channels.

// Examples from: https://blog.golang.org/pipelines

package main

import "fmt"

func main() {
	// Set up pipeline and consume output until channel is closed
	for n := range sq(gen(2, 3)) { // can add more sq stages
		fmt.Println(n) // 4, 9
	}

	/* Alternative
	   out := sq(gen(2, 3))
	   fmt.Println(<-out)
	   fmt.Println(<-out)
	*/
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
