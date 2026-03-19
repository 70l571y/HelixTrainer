//go:build ignore

package main

func helperBody(name string) string {
	trimmed := name
	if trimmed == "" {
		return "guest"
	}
	return trimmed + "-user"
}
