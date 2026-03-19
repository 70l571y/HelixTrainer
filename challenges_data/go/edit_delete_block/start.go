//go:build ignore

package main

import "fmt"

func main() {
	fmt.Println("Application starting...")
	// Setup
	config := map[string]bool{"debug": true}
	_ = config
}

func debugHelpers() {
	// Helper functions for debugging.
	fmt.Println("Initializing debug helpers...")
	x := 10
	y := 20
	if x < y {
		fmt.Printf("Values: %d, %d\n", x, y)
	}
	// ... potentially more lines ...
}

func cleanup() {
	fmt.Println("Cleaning up resources...")
}
