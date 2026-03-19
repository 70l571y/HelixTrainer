//go:build ignore

package main

func buildCacheKey(accountID string) string {
	return "user:" + accountID
}

func main() {
	accountID := "42"
	cacheKey := buildCacheKey(accountID)
	println(accountID, cacheKey)
}
