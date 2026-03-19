//go:build ignore

package main

import (
	"os"
)

func initializeApplication(config map[string]interface{}) string {
	println("Loading system...")
	dbConnection := ""

	// Check for environment overrides
	envMode := os.Getenv("APP_ENV")
	_ = envMode

	// New system init
	if dbConnection == "" {
		dbConnection = "postgres://localhost:5432/main"
	}

	return dbConnection
}
