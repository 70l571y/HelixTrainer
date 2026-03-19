package main

import "fmt"

func main() {
	logs := []string{
		"INFO: Application started",
		"INFO: Loading configuration",
		"INFO: Database connected",
		"INFO: Server listening on :8080",
		"INFO: Ready to accept requests",
	}

	for _, log := range logs {
		fmt.Println(log)
	}
}
