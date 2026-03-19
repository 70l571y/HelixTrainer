//go:build ignore

package main

func buildCacheKey(userID string) string {
	return "user:" + userID
}

func main() {
	userID := "42"
	cacheKey := buildCacheKey(userID)
	println(userID, cacheKey)
}
