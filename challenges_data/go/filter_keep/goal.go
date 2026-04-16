//go:build ignore

package main

import "fmt"

func processLogs() {
	logs := []string{
		"Error: Connection failed",
		"Error: Database timeout",
		"Error: Authentication failed",
	}

	for _, log := range logs {
		fmt.Println(log)
	}
}
