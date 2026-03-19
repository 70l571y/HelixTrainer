//go:build ignore

package main

// Copy the function 'parseLine' from parser_utils.go below this line
func parseLine(line string) (string, string) {
	if line == "" {
		return "", ""
	}
	// Simple parsing logic
	return "key", "value"
}

// Copy the struct 'ConfigLoader' from config_loader.go below this line
type ConfigLoader struct {
	filename string
	config   map[string]string
}

func NewConfigLoader(filename string) *ConfigLoader {
	return &ConfigLoader{
		filename: filename,
		config:   make(map[string]string),
	}
}
