//go:build ignore

package main

func loadState(ready bool) string {
	if !ready {
		return "fallback"
	}

	return "ready"
}
