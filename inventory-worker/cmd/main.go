package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Hello World")
		time.Sleep(2 * time.Second)
	}
}
