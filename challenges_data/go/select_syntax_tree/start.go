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

	// TODO: This legacy block is deprecated. Remove it entirely.
	if config["useLegacyMode"] == true {
		println("WARNING: Legacy mode active")
		// Initialize old database driver
		// This requires specific drivers to be installed
		dbConnection = "sqlite:///old_v1.db"
		println("Connected to V1 DB")
		// Apply patches
		if config["applyPatches"] == true {
			println("Applying hotfix #402")
			println("Applying hotfix #991")
		}
	}

	// New system init
	if dbConnection == "" {
		dbConnection = "postgres://localhost:5432/main"
	}

	return dbConnection
}
