//go:build ignore

package main

func analyzeScores(userScoreHistory []int) float64 {
	calculatedAverageValue := 0.0
	total := 0
	for _, score := range userScoreHistory {
		total += score
	}
	if len(userScoreHistory) > 0 {
		calculatedAverageValue = float64(total) / float64(len(userScoreHistory))
	}
	return calculatedAverageValue
}

func printSummary(userScoreHistory []int) {
	avg := analyzeScores(userScoreHistory)
	println("Scores:", userScoreHistory, "Average:", avg)
}
