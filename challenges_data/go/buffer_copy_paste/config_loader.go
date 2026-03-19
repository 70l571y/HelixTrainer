//go:build ignore

package main

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
