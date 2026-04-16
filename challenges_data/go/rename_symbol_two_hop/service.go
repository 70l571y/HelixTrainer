//go:build ignore

package main

func BuildBadge(scope string) string {
	label := "badge"
	return label + ":" + scope
}
