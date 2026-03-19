//go:build ignore

package main

func services() []string {
	return []string{
		"api",
		"worker",
		"billing",
	}
}
