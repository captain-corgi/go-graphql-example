package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello, World!")

	// A goroutine which cause goroutine leak
	//  this goroutines should cause race condition
	message := "Hello, World!"
	go func(s string) {
		time.Sleep(10 * time.Second)
		fmt.Println(s)
	}(message)

	time.Sleep(1 * time.Second)
	fmt.Println(message)
}
