package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	done := make(chan struct{})

	go func() {
		// Simulate some work
		time.Sleep(1 * time.Second)
		fmt.Println("Work completed")
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("Finished successfully")
	case <-ctx.Done():
		fmt.Println("Context cancelled or timed out:", ctx.Err())
	}
}
