//go:build ignore

package main

func parseLine(line string) (string, string) {
	if line == "" {
		return "", ""
	}
	// Simple parsing logic
	return "key", "value"
}
