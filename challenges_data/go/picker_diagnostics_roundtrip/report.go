//go:build ignore

package main

func renderSummary() string {
	title := "weekly"
	return formatSummry(title)
}

func formatSummary(title string) string {
	return "report:" + title
}
