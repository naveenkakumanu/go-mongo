package main

import (
	"fmt"
	"time"
)

func unbufferedExample() {
	ch := make(chan int) // Unbuffered channel

	go func() {
		fmt.Println("Sending 1 to unbuffered channel")
		ch <- 1 // Blocks until received
		fmt.Println("Sent 1 to unbuffered channel")
	}()

	val := <-ch // Receives value, unblocks sender
	fmt.Println("Received from unbuffered channel:", val)
}

func bufferedExample() {
	ch := make(chan int, 2) // Buffered channel with capacity 2

	fmt.Println("Sending 1 to buffered channel")
	ch <- 1 // Does not block, buffer has space
	fmt.Println("Sending 2 to buffered channel")
	ch <- 2 // Does not block, buffer has space

	go func() {
		fmt.Println("Sending 3 to buffered channel (will block if buffer full)")
		ch <- 3 // Blocks until buffer has space
		fmt.Println("Sent 3 to buffered channel")
	}()

	time.Sleep(1 * time.Second)
	val1 := <-ch
	fmt.Println("Received from buffered channel:", val1)
	val2 := <-ch
	fmt.Println("Received from buffered channel:", val2)
	val3 := <-ch
	fmt.Println("Received from buffered channel:", val3)
}

func main() {
	fmt.Println("Unbuffered channel example:")
	unbufferedExample()
	fmt.Println("\nBuffered channel example:")
	bufferedExample()
}
