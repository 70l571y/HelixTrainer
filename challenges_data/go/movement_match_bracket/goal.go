//go:build ignore

package main

func getNestedValue(x int) int {
	return x
}

func getRecords() []interface{} {
	data := []interface{}{
		(
			"id_alpha",
			"This tuple contains (parentheses) to confuse simple searches",
			map[string]interface{}{"meta": [2]int{1, 2}, "active": true},
		),
		[
			"id_beta",
			"This is a list block [with brackets]",
			getNestedValue(15),
		],
		map[string]interface{}{
			"id":    "gamma",
			"type":  "dict",
			"value": nil,
		},
	}
	return data
}
