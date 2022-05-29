# Timers


#### Repeat logic at an interval
```go
package main

import (
"context"
"fmt"
"time"
)

func main() {
    ctx := context.TODO()
	someDuration := 2 * time.Second
	
	timer := time.NewTimer(0)
	fmt.Print("Starting...")
	for {
		select {
		case <-timer.C:
		case <-ctx.Done():
			return
		}

		fmt.Printf("hey")
		timer.Reset(someDuration)
	}
}
```
We set timer to 0, so timer.C doesn't block and we execute the logic once immediately.

"Hey is printed", then timer restarts.
The select statement blocks while the timer is counting down.

- Q. How does this stop?


```go
package main

import (
"context"
"fmt"
"time"
)

func main() {
    ctx := context.TODO()
	someDuration := 2 * time.Second

	timer := time.NewTimer(0)
	defer timer.Stop()
	fmt.Print("Starting...")
	for {
		select {
		case <-timer.C:
		case <-ctx.Done():
			return
		}

		fmt.Printf("hey")
		timer.Reset(someDuration)
	}
}
```

- Q. When to add the defer?