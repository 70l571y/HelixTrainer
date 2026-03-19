//go:build ignore

package main

func services() []string {
	return []string{
		"svc:api",
		"svc:worker",
		"svc:billing",
	}
}
