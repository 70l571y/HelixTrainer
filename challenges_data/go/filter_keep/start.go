package main

import "fmt"

func processLogs() {
	logs := []string{
		"INFO: Starting application",
		"Error: Connection failed",
		"DEBUG: Loading config",
		"Error: Database timeout",
		"INFO: Retrying connection",
		"Error: Authentication failed",
		"INFO: Shutdown complete",
	}

	for _, log := range logs {
		fmt.Println(log)
	}
}
