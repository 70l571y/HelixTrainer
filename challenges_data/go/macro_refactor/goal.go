//go:build ignore

package main

import "fmt"

func main() {
	logs := []string{
		"log: INFO: Application started",
		"log: INFO: Loading configuration",
		"log: INFO: Database connected",
		"log: INFO: Server listening on :8080",
		"log: INFO: Ready to accept requests",
	}

	for _, log := range logs {
		fmt.Println(log)
	}
}
