//go:build ignore

package main

func analyzeScores(scores []int) float64 {
	average := 0.0
	total := 0
	for _, score := range scores {
		total += score
	}
	if len(scores) > 0 {
		average = float64(total) / float64(len(scores))
	}
	return average
}

func printSummary(scores []int) {
	avg := analyzeScores(scores)
	println("Scores:", scores, "Average:", avg)
}
