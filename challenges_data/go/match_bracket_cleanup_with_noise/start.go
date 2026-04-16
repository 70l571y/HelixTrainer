//go:build ignore

package main

func main() {
	tasks := []map[string]interface{}{
		{
			"id":      "alpha",
			"payload": []int{1, 2, 3}
			"note":    "payload uses [square brackets] in text too",
		},
		{
			"id": "beta",
			"meta": map[string]string{
				"mode": "fast",
				"kind": "audit",
			}
			"note": "this line mentions {braces} only as noise",
		},
	}

	_ = tasks
}
