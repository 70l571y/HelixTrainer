//go:build ignore

package main

func defaultConfig() map[string]interface{} {
	// These values are commented out but aligned
	userSettings := map[string]interface{}{
		"theme":       "dark",
		"notifications": true,
		"autoSave":    300,
	}
	return userSettings
}
