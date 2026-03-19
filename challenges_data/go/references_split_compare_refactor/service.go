//go:build ignore

package main

func runService() string {
	return normalizeName("service")
}

func normalizeName(name string) string {
	return name
}
