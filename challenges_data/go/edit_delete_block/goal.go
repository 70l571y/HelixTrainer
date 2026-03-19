//go:build ignore

package main

import "fmt"

func main() {
	fmt.Println("Application starting...")
	// Setup
	config := map[string]bool{"debug": true}
	_ = config
}

func cleanup() {
	fmt.Println("Cleaning up resources...")
}
