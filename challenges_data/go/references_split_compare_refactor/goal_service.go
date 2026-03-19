//go:build ignore

package main

func runService() string {
	return sanitizeName("service")
}

func sanitizeName(name string) string {
	return name
}
